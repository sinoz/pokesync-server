package entity

import "testing"

import "time"

func TestIntervalPolicy_Accumulation(t *testing.T) {
	policy := NewIntervalPolicy(250 * time.Millisecond)
	rate := 50 * time.Millisecond

	for i := 1; i <= 4; i++ {
		policy.Update(rate)

		expectedAccumulation := int64(i) * int64(rate)
		if int64(policy.accumulator) != expectedAccumulation {
			t.Errorf("expected accumulated time to equal %v ms", expectedAccumulation)
		}
	}

	policy.Update(rate)
	if policy.accumulator != 0*time.Millisecond {
		t.Error("expected accumulated time to equal 0 ms")
	}
}
