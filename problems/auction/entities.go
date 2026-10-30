package auction

import "time"

// Car is the item being auctioned. It's plain data: the auction refers to it,
// it doesn't own it (a car can be re-auctioned if it goes unsold).
type Car struct {
	ID    string
	Model string
	Year  int
}

// Dealer is a bidder. Bids store only the dealer's ID: a bid needs to know
// who placed it, not carry a copy of the dealer.
type Dealer struct {
	ID   string
	Name string
}

// Bid is owned by its Auction (composition), so it doesn't point back to it.
// Money is int64 rupees, never float: floats can't represent amounts exactly.
type Bid struct {
	DealerID string
	Amount   int64
	PlacedAt time.Time
}

// Result is what Close reports. Sold is false when nobody bid or the highest
// bid stayed below the reserve price.
type Result struct {
	Sold     bool
	WinnerID string
	Amount   int64
}
