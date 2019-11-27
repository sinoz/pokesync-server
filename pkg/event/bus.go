package event

import (
	"errors"
	"reflect"
)

// Receiver is the function that takes the parameters of an Event.
type Receiver reflect.Value

// Parameter is a parameter of an event passed to a Receiver.
type Parameter interface{}

// Topic is a type of Event topic to publish events to.
type Topic string

// Bus forwards events to their subscribers in a fire-and forget fashion.
type Bus interface {
	Subscribe(topic Topic, fn interface{}) error
	Publish(topic Topic, arguments ...Parameter)
	Unsubscribe(topic Topic, fn interface{})
	Collapse(topic Topic)
	Contains(topic Topic) bool
}

// SerialBus is an implementation of an event Bus that subscribes,
// unsubscribes and publishes events in a serial or sequential manner.
type SerialBus struct {
	receivers map[Topic][]Receiver
}

// NewSerialBus constructs a new instance of an event Bus that
// publishes messages serially or sequential.
func NewSerialBus() Bus {
	return &SerialBus{receivers: make(map[Topic][]Receiver)}
}

// Subscribe subscribes the given Receiver for the specified Topic.
func (bus *SerialBus) Subscribe(topic Topic, fn interface{}) error {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return errors.New("expected given value of 'fn' to equal type Func")
	}

	receivers := bus.receivers[topic]
	if receivers == nil {
		receivers = []Receiver{}
	}

	receivers = append(receivers, Receiver(reflect.ValueOf(fn)))
	bus.receivers[topic] = receivers

	return nil
}

// Publish publishes the given parameter list as an event to listeners
// of the specified Topic.
func (bus *SerialBus) Publish(topic Topic, params ...Parameter) {
	receivers := bus.receivers[topic]
	if receivers == nil {
		return
	}

	passedArguments := make([]reflect.Value, 0)
	for _, arg := range params {
		passedArguments = append(passedArguments, reflect.ValueOf(arg))
	}

	for _, rcv := range receivers {
		reflect.Value(rcv).Call(passedArguments)
	}
}

// Unsubscribe unsubscribes the given Receiver from the specified Topic.
func (bus *SerialBus) Unsubscribe(topic Topic, fn interface{}) {
	receivers := bus.receivers[topic]
	if receivers == nil {
		return
	}

	for i, rcv := range receivers {
		if rcv == Receiver(reflect.ValueOf(fn)) {
			receivers = append(receivers[:i], receivers[i+1:]...)
			break
		}
	}

	if len(receivers) == 0 {
		delete(bus.receivers, topic)
	} else {
		bus.receivers[topic] = receivers
	}
}

// Contains checks if there is a Topic registered of the specified kind.
func (bus *SerialBus) Contains(topic Topic) bool {
	return bus.receivers[topic] != nil
}

// Collapse collapses the given Receiver from the specified Topic.
func (bus *SerialBus) Collapse(topic Topic) {
	delete(bus.receivers, topic)
}
