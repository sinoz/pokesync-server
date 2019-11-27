package game

import "testing"

func TestClock_Pulse(t *testing.T) {
	clock := NewClock(60)
	clock.Pulse()
	if clock.seconds != 61 {
		t.Error("expected a second to have passed on the clock")
	}
}

func TestClock_PulseResetAfterMidnight(t *testing.T) {
	clock := NewClock(86401)
	clock.Pulse()
	if clock.seconds != 0 {
		t.Error("expected the clock to have reset its seconds")
	}
}

func TestClock_ReachedMidnightFalse(t *testing.T) {
	clock := NewClock(86399)
	if clock.ReachedMidnight() {
		t.Error("expected clock to have not reached midnight yet")
	}
}

func TestClock_ReachedMidnightTrue(t *testing.T) {
	clock := NewClock(86400)
	if !clock.ReachedMidnight() {
		t.Error("expected clock to have reached midnight")
	}
}

func TestClock_CurrentHour(t *testing.T) {
	clock := NewClock(60 * 60 * 14)
	hour := clock.CurrentHour()
	if hour != 14 {
		t.Errorf("expected hour to equal %v but was %v instead", 14, hour)
	}
}

func TestClock_CurrentMinute(t *testing.T) {
	clock := NewClock(60 * 24)
	minute := clock.CurrentMinute()
	if minute != 24 {
		t.Errorf("expected minute to equal %v but was %v instead", 24, minute)
	}
}
