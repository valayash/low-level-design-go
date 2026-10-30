package auction

// Notifier is the only thing the auction knows about notifications. SMS,
// email or push each implement it; the auction never changes when a channel
// is added (Dependency Inversion + Open/Closed).
type Notifier interface {
	Send(userID, msg string) error
}

// NoopNotifier discards messages. Useful when nobody needs to be told.
type NoopNotifier struct{}

func (NoopNotifier) Send(string, string) error { return nil }
