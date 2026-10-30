package auction

import (
	"fmt"
	"sync"
	"time"
)

// Auction owns its bids and its lifecycle. Every method that reads or changes
// state takes the mutex, so two dealers bidding at the same instant can never
// both become the highest bidder.
type Auction struct {
	mu sync.Mutex

	id            string
	car           *Car
	startingPrice int64
	reservePrice  int64
	minIncrement  int64
	duration      time.Duration
	endTime       time.Time

	status  Status
	bids    []Bid
	highest *Bid

	notifier Notifier
	now      func() time.Time // injectable clock so tests can control time
}

// Config groups NewAuction's inputs so the constructor doesn't take eight
// positional arguments.
type Config struct {
	ID            string
	Car           *Car
	StartingPrice int64
	ReservePrice  int64
	MinIncrement  int64
	Duration      time.Duration
	Notifier      Notifier         // optional; defaults to NoopNotifier
	Now           func() time.Time // optional; defaults to time.Now
}

// NewAuction validates its input so an invalid auction can't exist, and
// starts in Scheduled.
func NewAuction(c Config) (*Auction, error) {
	if c.Car == nil {
		return nil, ErrNilCar
	}
	if c.StartingPrice <= 0 || c.StartingPrice > c.ReservePrice {
		return nil, ErrInvalidPrice
	}
	if c.MinIncrement <= 0 {
		return nil, ErrInvalidIncrement
	}
	if c.Duration <= 0 {
		return nil, ErrInvalidDuration
	}
	if c.Notifier == nil {
		c.Notifier = NoopNotifier{}
	}
	if c.Now == nil {
		c.Now = time.Now
	}

	return &Auction{
		id:            c.ID,
		car:           c.Car,
		startingPrice: c.StartingPrice,
		reservePrice:  c.ReservePrice,
		minIncrement:  c.MinIncrement,
		duration:      c.Duration,
		status:        Scheduled,
		notifier:      c.Notifier,
		now:           c.Now,
	}, nil
}

// transition is the single place status changes; callers must hold mu.
func (a *Auction) transition(to Status) error {
	if !canMove(a.status, to) {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidTransition, a.status, to)
	}
	a.status = to
	return nil
}

// Start opens the auction for bids; the clock starts now.
func (a *Auction) Start() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if err := a.transition(Live); err != nil {
		return err
	}
	a.endTime = a.now().Add(a.duration)
	return nil
}

// PlaceBid applies the bid rules in a fixed order: the auction must be live
// and not past its end time, the dealer must not outbid themselves, and the
// amount must reach the starting price (first bid) or beat the highest by the
// minimum increment.
//
// The outbid notification is sent after the lock is released, so a slow
// notifier never blocks other bidders.
func (a *Auction) PlaceBid(dealerID string, amount int64) error {
	a.mu.Lock()

	if a.status != Live {
		a.mu.Unlock()
		return ErrNotLive
	}
	now := a.now()
	if !now.Before(a.endTime) {
		a.mu.Unlock()
		return ErrAuctionEnded
	}

	minAmount := a.startingPrice
	if a.highest != nil {
		if a.highest.DealerID == dealerID {
			a.mu.Unlock()
			return ErrAlreadyHighest
		}
		minAmount = a.highest.Amount + a.minIncrement
	}
	if amount < minAmount {
		a.mu.Unlock()
		return fmt.Errorf("%w: minimum is %d", ErrBidTooLow, minAmount)
	}

	var outbid string
	if a.highest != nil {
		outbid = a.highest.DealerID
	}
	bid := Bid{DealerID: dealerID, Amount: amount, PlacedAt: now}
	a.bids = append(a.bids, bid)
	a.highest = &bid // own copy: a pointer into a.bids would go stale when append reallocates
	a.mu.Unlock()

	if outbid != "" {
		_ = a.notifier.Send(outbid, fmt.Sprintf("You were outbid on auction %s: new highest is %d", a.id, amount))
	}
	return nil
}

// Close ends the auction. The car sells only if the highest bid reached the
// reserve price; the winner is notified after the lock is released.
func (a *Auction) Close() (Result, error) {
	a.mu.Lock()

	if err := a.transition(Closed); err != nil {
		a.mu.Unlock()
		return Result{}, err
	}

	var res Result
	if a.highest != nil && a.highest.Amount >= a.reservePrice {
		res = Result{Sold: true, WinnerID: a.highest.DealerID, Amount: a.highest.Amount}
	}
	a.mu.Unlock()

	if res.Sold {
		_ = a.notifier.Send(res.WinnerID, fmt.Sprintf("You won auction %s at %d", a.id, res.Amount))
	}
	return res, nil
}

// Cancel stops a scheduled or live auction. Bids already placed are kept for audit.
func (a *Auction) Cancel() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.transition(Cancelled)
}

// Status reports the current lifecycle state.
func (a *Auction) Status() Status {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.status
}

// Highest returns a copy of the current highest bid, and false if there is none.
// Returning a copy means callers can't change the auction's bid from outside.
func (a *Auction) Highest() (Bid, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.highest == nil {
		return Bid{}, false
	}
	return *a.highest, true
}
