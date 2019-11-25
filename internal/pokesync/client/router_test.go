package client

import (
	"testing"
	"time"
)

func TestRouter_Subscribe(t *testing.T) {
	router := NewRouter()
	router.Subscribe("hello")

	if router.TopicCount() != 1 {
		t.Error("expected topic count to equal 1")
	}
}

func TestRouter_SubscribeMailboxToTopic(t *testing.T) {
	router := NewRouter()

	m := router.Subscribe("hello")

	for i := 0; i < 25; i++ {
		err := router.SubscribeMailboxToTopic("hello", m)
		if err != nil {
			t.Fatal(err)
		}
	}

	if router.TopicCount() != 1 {
		t.Error("expected topic count to equal 1")
	}

	if router.SizeOf("hello") != 26 {
		t.Errorf("expected subscriber count of topic %v to equal %v\n", "hello", 26)
	}
}

func TestRouter_Publish(t *testing.T) {
	router := NewRouter()

	m1 := router.Subscribe("hello")
	m2 := router.Subscribe("hello1")

	go func() {
		router.Publish("hello", Mail{Payload: 1})
	}()

	go func() {
		router.Publish("hello1", Mail{Payload: 2})
	}()

	time.Sleep(time.Second)

	if x := <-m1; x.Payload != 1 {
		t.Error("expected message value from mailbox to equal 1")
	}

	if x := <-m2; x.Payload != 2 {
		t.Error("expected message value from mailbox to equal 1")
	}
}

func TestRouter_Collapse(t *testing.T) {
	router := NewRouter()

	router.Subscribe("hello")
	router.Subscribe("hello1")

	if router.TopicCount() != 2 {
		t.Error("expected topic count to equal 2")
	}

	router.Collapse("hello")
	if router.TopicCount() != 1 {
		t.Error("expected topic count to equal 1")
	}

	router.Collapse("hello1")
	if router.TopicCount() != 0 {
		t.Error("expected topic count to equal 0")
	}
}

func TestRouter_Unsubscribe(t *testing.T) {
	router := NewRouter()

	m1 := router.Subscribe("hello")
	m2 := router.Subscribe("hello1")

	if router.TopicCount() != 2 {
		t.Error("expected topic count to equal 2")
	}

	router.Unsubscribe("hello", m1)
	if router.TopicCount() != 1 {
		t.Error("expected topic count to equal 1")
	}

	router.Unsubscribe("hello1", m2)
	if router.TopicCount() != 0 {
		t.Error("expected topic count to equal 0")
	}
}
