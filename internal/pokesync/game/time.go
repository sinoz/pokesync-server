package game

import "time"

// The amount of minutes in an hour.
const MinutesInAnHour = 60

// The amount of seconds in a minute.
const SecondsInAMinute = 60

// The amount of hours in a day.
const HoursInADay = 24

// The amount of seconds in a day.
const SecondsInADay = SecondsInAMinute * MinutesInAnHour * HoursInADay

// Clock points to a specific hour and minute on a clock.
type Clock struct {
	seconds int
}

// Synchronizer synchronizes a Clock with an external standard uniform
// of time, such as an actual time zone.
type Synchronizer interface {
	Synchronize() (*Clock, error)
}

// GMT0Synchronizer synchronizes with the GMT0 timezone.
type GMT0Synchronizer struct{}

// NewClock constructs a new Clock.
func NewClock(seconds int) *Clock {
	return &Clock{seconds: seconds}
}

// Pulse is called every game pulse.
func (clock *Clock) Pulse() {
	clock.incrementSecond()
	if clock.reachedMidnight() {
		clock.resetSeconds()
	}
}

// incrementSecond adds a second to the 'seconds' field in Clock.
func (clock *Clock) incrementSecond() {
	clock.seconds++
}

// resetSeconds resets the 'seconds' field in Clock back to zero.
func (clock *Clock) resetSeconds() {
	clock.seconds = 0
}

// reachedMidnight checks if the Clock has reached midnight.
func (clock Clock) reachedMidnight() bool {
	return clock.seconds >= SecondsInADay
}

// Returns the current hour in a day.
func (clock Clock) currentHour() int {
	return (clock.seconds / (MinutesInAnHour * SecondsInAMinute)) % HoursInADay
}

// Returns the current minute in an hour.
func (clock Clock) currentMinute() int {
	return (clock.seconds / SecondsInAMinute) % SecondsInAMinute
}

// NewGMT0Synchronizer constructs a new Synchronizer that synchronizes
// with the GMT+0 timezone.
func NewGMT0Synchronizer() Synchronizer {
	return new(GMT0Synchronizer)
}

// Synchronize synchronizes with the GMT0 timezone, returning a Clock
// with the amount of seconds that has passed since midnight.
func (gmt *GMT0Synchronizer) Synchronize() (*Clock, error) {
	location, err := time.LoadLocation("GMT-0")
	if err != nil {
		return nil, err
	}

	t := time.Now().In(location)

	hours := t.Hour()
	minutes := t.Minute()
	seconds := t.Second()

	return NewClock((hours * MinutesInAnHour * SecondsInADay) + (minutes * SecondsInAMinute) + seconds), nil
}
