package game

import (
	"time"

	"gitlab.com/pokesync/game-service/internal/game-service/game/entity"
	"gitlab.com/pokesync/game-service/internal/game-service/game/transport"
)

// MinutesInAnHour is the amount of minutes in an hour.
const MinutesInAnHour = 60

// SecondsInAMinute is the amount of seconds in a minute.
const SecondsInAMinute = 60

// HoursInADay is the amount of hours in a day.
const HoursInADay = 24

// SecondsInADay is the amount of seconds in a day.
const SecondsInADay = SecondsInAMinute * MinutesInAnHour * HoursInADay

// Clock points to a specific hour and minute on a clock.
type Clock struct {
	seconds int
}

// ClockSynchronizer synchronizes a Clock with an external standard uniform
// of time, such as an actual time zone.
type ClockSynchronizer interface {
	Synchronize() (*Clock, error)
}

// GMT0Synchronizer synchronizes with the GMT0 timezone.
type GMT0Synchronizer struct{}

// DayNightProcessor processes the transition of the day-night time periods.
type DayNightProcessor struct {
	clock        *Clock
	synchronizer ClockSynchronizer

	lastMinute int
}

// NewClock constructs a new Clock.
func NewClock(seconds int) *Clock {
	return &Clock{seconds: seconds}
}

// NewDayNightSystem constructs a System that processes the transitions
// between day-and night. The system processes time four times as fast
// as real-time does. This means that there are four transitions between
// day and night, a day.
func NewDayNightSystem(clockRate time.Duration, synchronizer ClockSynchronizer) *entity.System {
	return entity.NewSystem(entity.NewIntervalPolicy(clockRate), NewDayNightProcessor(synchronizer))
}

// NewDayNightProcessor processes the day-and night transitions.
func NewDayNightProcessor(synchronizer ClockSynchronizer) *DayNightProcessor {
	return &DayNightProcessor{synchronizer: synchronizer}
}

// Pulse is called every game pulse.
func (clock *Clock) Pulse() {
	clock.incrementSecond()
	if clock.ReachedMidnight() {
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

// ReachedMidnight checks if the Clock has reached midnight.
func (clock Clock) ReachedMidnight() bool {
	return clock.seconds >= SecondsInADay
}

// CurrentHour returns the current hour in a day.
func (clock Clock) CurrentHour() int {
	return (clock.seconds / (MinutesInAnHour * SecondsInAMinute)) % HoursInADay
}

// CurrentMinute returns the current minute in an hour.
func (clock Clock) CurrentMinute() int {
	return (clock.seconds / SecondsInAMinute) % SecondsInAMinute
}

// NewGMT0Synchronizer constructs a new ClockSynchronizer that synchronizes
// with the GMT+0 timezone.
func NewGMT0Synchronizer() ClockSynchronizer {
	return new(GMT0Synchronizer)
}

// Synchronize synchronizes with the GMT0 timezone, returning a Clock
// with the amount of seconds that has passed since midnight.
func (gmt *GMT0Synchronizer) Synchronize() (*Clock, error) {
	location, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, err
	}

	t := time.Now().In(location)

	hours := t.Hour()
	minutes := t.Minute()
	seconds := t.Second()

	return NewClock((hours * MinutesInAnHour * SecondsInADay) + (minutes * SecondsInAMinute) + seconds), nil
}

// AddedToWorld is called when the System of this Processor is added
// to the game World.
func (processor *DayNightProcessor) AddedToWorld(world *entity.World) (err error) {
	processor.clock, err = processor.synchronizer.Synchronize()
	return err
}

// RemovedFromWorld is called when the System of this Processor is removed
// from the game World.
func (processor *DayNightProcessor) RemovedFromWorld(world *entity.World) error {
	return nil
}

// Update is called every game pulse to check if entities need their map view
// refreshed and if so, refreshes them.
func (processor *DayNightProcessor) Update(world *entity.World, deltaTime time.Duration) error {
	processor.clock.Pulse()

	entities := world.GetEntitiesFor(processor)
	for _, ent := range entities {
		sessionComponent := ent.GetComponent(SessionTag).(*SessionComponent)

		currentHour := processor.clock.CurrentHour()
		currentMinute := processor.clock.CurrentMinute()
		if currentMinute != processor.lastMinute {
			sessionComponent.session.QueueEvent(&transport.SetServerTime{
				Hour:   byte(currentHour),
				Minute: byte(currentMinute),
			})

			processor.lastMinute = currentMinute
		}
	}

	return nil
}

// Components returns a pack of ComponentTag's the DayNightProcessor has
// interest in.
func (processor *DayNightProcessor) Components() entity.ComponentTag {
	return WaryOfTimeTag | SessionTag
}
