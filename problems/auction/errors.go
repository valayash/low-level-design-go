package auction

import "errors"

// Sentinel errors so callers can tell why a call failed with errors.Is.
var (
	ErrInvalidPrice      = errors.New("starting price must be positive and not above the reserve price")
	ErrInvalidIncrement  = errors.New("minimum increment must be positive")
	ErrInvalidDuration   = errors.New("duration must be positive")
	ErrNilCar            = errors.New("car must not be nil")
	ErrInvalidTransition = errors.New("invalid status transition")
	ErrNotLive           = errors.New("auction is not live")
	ErrAuctionEnded      = errors.New("auction end time has passed")
	ErrBidTooLow         = errors.New("bid is too low")
	ErrAlreadyHighest    = errors.New("dealer is already the highest bidder")
)
