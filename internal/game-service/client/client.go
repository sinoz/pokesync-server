package client

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"reflect"

	"github.com/google/uuid"
)

// Config holds configurations specific for Client's.
type Config struct {
	MessageCodec Codec

	ReadBufferSize  int
	WriteBufferSize int

	CommandLimit int
}

// ID is the unique identifier of the Client.
type ID uuid.UUID

// Client represents a connected game client.
type Client struct {
	config Config

	ID ID

	connection net.Conn

	reader *bufio.Reader
	writer *bufio.Writer

	commands chan command

	codec Codec
}

// flushCommand is a cached instance of the 'flush' type.
var (
	flushCommand     command = flush{}
	terminateCommand command = terminate{}
)

// command is a client command.
type command interface{}

// send is a command to marshal and queue a message.
type send struct {
	message Message
}

// terminate is a command to terminate the socket connection with
// the client.
type terminate struct{}

// flush is a command to flush queued up byte contents to the socket
// connection.
type flush struct{}

// TerminationTopic is the topic of the Terminated event.
const TerminationTopic = "client_terminated"

// Terminated is an event that is broadcasted to all services to notify
// them of a Client's termination.
type Terminated struct {
	ID ID
}

// BuildNumber is the build number of the game client.
type BuildNumber int

// NewClient constructs a new Client for the given connection.
func NewClient(connection net.Conn, config Config) *Client {
	return &Client{
		ID:         ID(uuid.New()),
		connection: connection,

		reader: bufio.NewReaderSize(connection, config.ReadBufferSize),
		writer: bufio.NewWriterSize(connection, config.WriteBufferSize),

		commands: make(chan command, config.CommandLimit),

		config: config,
		codec:  config.MessageCodec,
	}
}

// Pull pulls a single message from the underlying socket connection and
// publishes the message to the given Router. Returns an error if:
//   - There was an issue with demarshalling a message,
//   - The underlying socket connection is closed,
//   - Or if the Router was unable to publish the message.
//
// The Context that is given as a parameter is passed on to every message.
// This Context is of importance for services to be notified of when the
// Client's underlying connection has been terminated. This allows services
// to interrupt any remaining operation that is associated with the
// disconnected Client.
func (c *Client) Pull(ctx context.Context, cancel context.CancelFunc, router *Router) error {
	packet, err := ForkPacket(c.reader)
	if err != nil {
		// notify all services that possess the context that they should
		// abandon the associated request they are processing as the client
		// has disconnected.
		cancel()

		// and also notify these services of this termination, should they
		// already have established a session for this terminated Client.
		router.Publish(TerminationTopic, Mail{
			Client:  c,
			Context: ctx,
			Payload: Terminated{ID: c.ID},
		})

		return err
	}

	config, exists := c.codec.GetConfig(packet.Kind)
	if !exists {
		return fmt.Errorf("no MessageConfig associated with packet of kind %v", packet.Kind)
	}

	message := config.New()
	message.Demarshal(packet)

	delivered := router.Publish(config.Topic, Mail{
		Client:  c,
		Context: ctx,
		Payload: message,
	})

	if !delivered {
		return fmt.Errorf("failed to deliver message of type %v with topic %v to a recipient", reflect.TypeOf(message), config.Topic)
	}

	return nil
}

// Push handles a single client push message, which may involve writing bytes
// to the underlying socket connection. Push may return an error if:
//   - The socket connection is closed,
//	 - Or if there was an issue with marshalling a message,
func (c *Client) Push(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil

	case command := <-c.commands:
		switch cmd := command.(type) {
		case send:
			config := cmd.message.GetConfig()
			if _, exists := c.codec.GetConfig(config.Kind); !exists {
				return fmt.Errorf("no MessageConfig associated with packet of kind %v", config.Kind)
			}

			packet := &Packet{Kind: config.Kind, Bytes: cmd.message.Marshal()}
			if err := JoinPacket(c.writer, packet); err != nil {
				// disconnected. when the writer reaches its capacity, it
				// automatically performs a flush to be able to write more
				// bytes, which may produce an error if the socket connection
				// is dropped
				return err
			}

		case flush:
			if err := c.writer.Flush(); err != nil {
				return err
			}

		case terminate:
			if err := c.connection.Close(); err != nil {
				return err
			}

		default:
			return fmt.Errorf("unexpected push command of type %v", reflect.TypeOf(cmd))
		}
	}

	return nil
}

// Send calls for the given Message to be marshalled and sent across the wire
// when a call for a flush occurs.
func (c *Client) Send(message Message) {
	c.commands <- send{message: message}
}

// SendNow calls for the given Message to be marshalled and sent directly
// across the wire by triggering a flush call.
func (c *Client) SendNow(message Message) {
	c.Send(message)
	c.Flush()
}

// Terminate calls for a termination of the client.
func (c *Client) Terminate() {
	c.commands <- terminateCommand
	close(c.commands)
}

// Flush calls for a flush of queued up bytes.
func (c *Client) Flush() {
	c.commands <- flushCommand
}

// IsUpToDateWith returns whether this BuildNumber is up-to-date with the
// given BuildNumber value.
func (b BuildNumber) IsUpToDateWith(other BuildNumber) bool {
	return b >= other
}
