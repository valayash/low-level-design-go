package auction

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

// recorder is a fake Notifier that remembers what was sent.
type recorder struct {
	mu   sync.Mutex
	sent map[string][]string
}

func (r *recorder) Send(userID, msg string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.sent == nil {
		r.sent = map[string][]string{}
	}
	r.sent[userID] = append(r.sent[userID], msg)
	return nil
}

func (r *recorder) count(userID string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.sent[userID])
}

// fakeClock lets tests move time forward without sleeping.
type fakeClock struct{ t time.Time }

func (c *fakeClock) now() time.Time { return c.t }

func newLiveAuction(t *testing.T) (*Auction, *recorder, *fakeClock) {
	t.Helper()
	rec := &recorder{}
	clk := &fakeClock{t: time.Date(2026, 10, 30, 18, 0, 0, 0, time.UTC)}
	a, err := NewAuction(Config{
		ID:            "A1",
		Car:           &Car{ID: "C1", Model: "Swift", Year: 2019},
		StartingPrice: 3_00_000,
		ReservePrice:  4_00_000,
		MinIncrement:  1_000,
		Duration:      30 * time.Minute,
		Notifier:      rec,
		Now:           clk.now,
	})
	if err != nil {
		t.Fatalf("NewAuction: %v", err)
	}
	if err := a.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	return a, rec, clk
}

func TestNewAuctionValidates(t *testing.T) {
	car := &Car{ID: "C1"}
	cases := []struct {
		name string
		cfg  Config
		want error
	}{
		{"nil car", Config{StartingPrice: 1, ReservePrice: 1, MinIncrement: 1, Duration: time.Minute}, ErrNilCar},
		{"start above reserve", Config{Car: car, StartingPrice: 5, ReservePrice: 4, MinIncrement: 1, Duration: time.Minute}, ErrInvalidPrice},
		{"zero increment", Config{Car: car, StartingPrice: 1, ReservePrice: 1, Duration: time.Minute}, ErrInvalidIncrement},
		{"zero duration", Config{Car: car, StartingPrice: 1, ReservePrice: 1, MinIncrement: 1}, ErrInvalidDuration},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := NewAuction(tc.cfg); !errors.Is(err, tc.want) {
				t.Fatalf("got %v, want %v", err, tc.want)
			}
		})
	}
}

func TestBidRules(t *testing.T) {
	a, rec, _ := newLiveAuction(t)

	if err := a.PlaceBid("d1", 2_99_999); !errors.Is(err, ErrBidTooLow) {
		t.Fatalf("below starting price: got %v", err)
	}
	if err := a.PlaceBid("d1", 3_00_000); err != nil {
		t.Fatalf("first bid at starting price: %v", err)
	}
	if err := a.PlaceBid("d1", 3_50_000); !errors.Is(err, ErrAlreadyHighest) {
		t.Fatalf("outbidding yourself: got %v", err)
	}
	if err := a.PlaceBid("d2", 3_00_999); !errors.Is(err, ErrBidTooLow) {
		t.Fatalf("below increment: got %v", err)
	}
	if err := a.PlaceBid("d2", 3_01_000); err != nil {
		t.Fatalf("valid outbid: %v", err)
	}
	if rec.count("d1") != 1 {
		t.Fatalf("d1 should be told they were outbid once, got %d", rec.count("d1"))
	}
	if h, _ := a.Highest(); h.DealerID != "d2" || h.Amount != 3_01_000 {
		t.Fatalf("highest = %+v", h)
	}
}

func TestBidsRejectedOutsideLive(t *testing.T) {
	a, err := NewAuction(Config{Car: &Car{ID: "C1"}, StartingPrice: 1, ReservePrice: 1, MinIncrement: 1, Duration: time.Minute})
	if err != nil {
		t.Fatal(err)
	}
	if err := a.PlaceBid("d1", 10); !errors.Is(err, ErrNotLive) {
		t.Fatalf("bid on scheduled auction: got %v", err)
	}

	live, _, clk := newLiveAuction(t)
	clk.t = clk.t.Add(30 * time.Minute)
	if err := live.PlaceBid("d1", 3_00_000); !errors.Is(err, ErrAuctionEnded) {
		t.Fatalf("bid after end time: got %v", err)
	}
}

func TestCloseSoldAndUnsold(t *testing.T) {
	a, rec, _ := newLiveAuction(t)
	_ = a.PlaceBid("d1", 3_00_000)
	_ = a.PlaceBid("d2", 4_00_000)

	res, err := a.Close()
	if err != nil {
		t.Fatal(err)
	}
	if !res.Sold || res.WinnerID != "d2" || res.Amount != 4_00_000 {
		t.Fatalf("result = %+v", res)
	}
	if rec.count("d2") != 1 {
		t.Fatal("winner should be notified")
	}

	b, _, _ := newLiveAuction(t)
	_ = b.PlaceBid("d1", 3_50_000) // below the 4L reserve
	if res, _ := b.Close(); res.Sold {
		t.Fatalf("below reserve must not sell: %+v", res)
	}
}

func TestInvalidTransitions(t *testing.T) {
	a, _, _ := newLiveAuction(t)
	if _, err := a.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := a.Close(); !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("closing twice: got %v", err)
	}
	if err := a.Cancel(); !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("cancelling a closed auction: got %v", err)
	}
	if err := a.Start(); !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("restarting a closed auction: got %v", err)
	}
}

// Many dealers bid the same amount at once: exactly one may win.
func TestConcurrentBidsOneWinner(t *testing.T) {
	a, _, _ := newLiveAuction(t)

	var wg sync.WaitGroup
	var mu sync.Mutex
	accepted := 0
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if a.PlaceBid(fmt.Sprintf("d%d", i), 3_00_000) == nil {
				mu.Lock()
				accepted++
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()

	if accepted != 1 {
		t.Fatalf("expected exactly one accepted bid, got %d", accepted)
	}
}
