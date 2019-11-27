package event

import (
	"testing"
)

func TestSerialBus_Subscribe(t *testing.T) {
	bus := NewSerialBus()
	bus.Subscribe("hello:world", func(y int) {
	})

	if !bus.Contains("hello:world") {
		t.Error("no topic going by 'hello:world' registered")
	}
}

func TestSerialBus_Publish(t *testing.T) {
	var x int
	cb := func(y int) {
		x += y
	}

	bus := NewSerialBus()
	bus.Subscribe("hello:world", cb)
	bus.Publish("hello:world", 5)

	if x != 5 {
		t.Errorf("expected x to have a value equal of %v but was %v instead", 5, x)
	}
}

func TestSerialBus_Unsubscribe(t *testing.T) {
	cb := func(y int) {
	}

	bus := NewSerialBus()
	bus.Subscribe("hello:world", cb)

	if !bus.Contains("hello:world") {
		t.Error("no topic going by 'hello:world' registered")
	}

	bus.Unsubscribe("hello:world", cb)
	if bus.Contains("hello:world") {
		t.Error("expected topic going by 'hello:world' to not exist")
	}
}
