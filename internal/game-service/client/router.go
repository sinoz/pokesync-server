package client

import (
	"context"
	"sync"
	"time"
)

// Topic is a type of Message
type Topic string

// Payload is the carried payload or the underlying message
// that is being sent from one point to another.
type Payload interface{}

// Mail is a message wrapping envelope that also holds the return address
// and/or sender of the underlying message (or payload).
type Mail struct {
	Client  *Client
	Context context.Context
	Payload Payload
}

// Mailbox is a channel of Message's that each subscriber gets.
type Mailbox chan Mail

// RouterConfig holds configurations specific for Router's.
type RouterConfig struct {
	PublicationTimeout time.Duration
}

// Router enroutes messages to their subscribers' mailboxes.
type Router struct {
	config    RouterConfig
	mailboxes map[Topic][]Mailbox
	mutex     *sync.RWMutex
}

// NewRouter constructs a new instance of Router.
func NewRouter(config RouterConfig) *Router {
	return &Router{
		config:    config,
		mailboxes: make(map[Topic][]Mailbox),
		mutex:     &sync.RWMutex{},
	}
}

// CreateMailbox constructs a Mailbox of messages.
func (r *Router) CreateMailbox() Mailbox {
	return make(Mailbox)
}

// Subscribe subscribes a new Mailbox to the specified topic, which is
// then returned so that the client caller can make use of the mailbox.
func (r *Router) Subscribe(topic Topic) Mailbox {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.mailboxes[topic] == nil {
		r.mailboxes[topic] = []Mailbox{}
	}

	mailbox := r.CreateMailbox()
	r.mailboxes[topic] = append(r.mailboxes[topic], mailbox)

	return mailbox
}

// SubscribeMailboxToTopic subscribes the given mailbox to an existing topic.
func (r *Router) SubscribeMailboxToTopic(topic Topic, mailbox Mailbox) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.mailboxes[topic] == nil {
		r.mailboxes[topic] = []Mailbox{}
	}

	r.mailboxes[topic] = append(r.mailboxes[topic], mailbox)
}

// Unsubscribe unsubscribes the given Mailbox from the specified topic.
func (r *Router) Unsubscribe(topic Topic, mailbox Mailbox) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	mailboxes := r.mailboxes[topic]
	for index, m := range mailboxes {
		if m == mailbox {
			r.mailboxes[topic] = append(mailboxes[:index], mailboxes[index+1:]...)
			close(m)

			break
		}
	}

	if len(r.mailboxes[topic]) == 0 {
		delete(r.mailboxes, topic)
	}
}

// Collapse collapses the specified topic, closing all subscriber mailboxes.
func (r *Router) Collapse(topic Topic) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, mailbox := range r.mailboxes[topic] {
		close(mailbox)
	}

	delete(r.mailboxes, topic)
}

// TopicCount returns the amount of topics that subscribers are subscribed to.
func (r *Router) TopicCount() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return len(r.mailboxes)
}

// SizeOf returns the amount of subscribers that are subscribed to the
// specified topic.
func (r *Router) SizeOf(topic Topic) int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return len(r.mailboxes[topic])
}

// Broadcast publishes the given Mail to all registered recipients in
// the router, regardless if they are interested in it or not.
func (r *Router) Broadcast(mail Mail) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, recipients := range r.mailboxes {
		for _, mailbox := range recipients {
			select {
			case mailbox <- mail:
				// mail delivered to recipient
			case <-time.After(r.config.PublicationTimeout):
				// skip this recipient then
				continue
			}
		}
	}
}

// Publish publishes the given Message on the specified Topic, which may
// broadcast the message across multiple channels. Returns whether the
// mail was delivered to a recipient or not.
func (r *Router) Publish(topic Topic, mail Mail) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	mailboxes := r.mailboxes[topic]
	for _, mailbox := range mailboxes {
		select {
		case mailbox <- mail:
			// mail delivered to recipient
		case <-time.After(r.config.PublicationTimeout):
			return false
		}
	}

	return true
}
