package auction

// Status is the auction's lifecycle (State pattern, enum + transition map).
type Status int

const (
	Scheduled Status = iota
	Live
	Closed
	Cancelled
)

func (s Status) String() string {
	switch s {
	case Scheduled:
		return "scheduled"
	case Live:
		return "live"
	case Closed:
		return "closed"
	case Cancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}

// allowed lists every legal move. Closed and Cancelled are final, so they have
// no entries. Adding a state (say Paused) is a change to this map only.
var allowed = map[Status][]Status{
	Scheduled: {Live, Cancelled},
	Live:      {Closed, Cancelled},
}

func canMove(from, to Status) bool {
	for _, s := range allowed[from] {
		if s == to {
			return true
		}
	}
	return false
}
