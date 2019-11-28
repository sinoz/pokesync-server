package status

// Notifier notifies external services of this server's online status.
type Notifier interface {
	Notify(parameters Parameters) error
}

// VoidNotifier is a type of Notifier that doesn't do anything with
// the status Parameters it is given.
type VoidNotifier struct{}

// RedisNotifier is a Notifier that stores status information in a
// Redis in-memory server instance.
type RedisNotifier struct {
	// TODO
}

// NewVoidNotifier constructs a new instance of a VoidNotifier.
func NewVoidNotifier() *VoidNotifier {
	return &VoidNotifier{}
}

// NewRedisNotifier constructs a new instance of a RedisNotifier.
func NewRedisNotifier() *RedisNotifier {
	return &RedisNotifier{}
}

// Notify doesn't do anything as intended.
func (notifier *VoidNotifier) Notify(parameters Parameters) error {
	return nil
}

// Notify stores the given set of Parameters into the Redis store. May
// return an error if something went wrong.
func (notifier *RedisNotifier) Notify(parameters Parameters) error {
	return nil
}
