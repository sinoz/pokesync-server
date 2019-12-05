package client

import (
	"bufio"
	"context"
	"net"
	"reflect"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Config holds configurations specific for Client's.
type Config struct {
	MessageCodec Codec
	Log          *zap.SugaredLogger

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

	router *Router
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
func NewClient(connection net.Conn, config Config, router *Router) *Client {
	return &Client{
		connection: connection,

		reader: bufio.NewReaderSize(connection, config.ReadBufferSize),
		writer: bufio.NewWriterSize(connection, config.WriteBufferSize),

		commands: make(chan command, config.CommandLimit),

		config: config,
		codec:  config.MessageCodec,

		router: router,
	}
}

// Pull pulls messages from the underlying socket connection until the socket
// connection is closed. The given Context is passed on to every message for
// services to be notified of when the Client's underlying connection has been
// terminated. This allows services to interrupt any remaining operation that
// is associated with the disconnected Client.
func (c *Client) Pull(ctx context.Context, cancel context.CancelFunc) {
	for {
		packet, err := ForkPacket(c.reader)
		if err != nil {
			// notify all services that possess the context that they should
			// abandon the associated request they are processing as the client
			// has disconnected.
			cancel()

			return
		}

		config, exists := c.codec.GetConfig(packet.Kind)
		if !exists {
			c.config.Log.Errorf("No MessageConfig associated with packet of kind %v\n", packet.Kind)
			continue
		}

		message := config.New()
		message.Demarshal(packet)

		mail := Mail{Client: c, Context: ctx, Payload: message}
		delivered := c.router.Publish(config.Topic, mail)
		if !delivered {
			c.config.Log.Errorf("Failed to deliver message of type %v with topic %v to a recipient", reflect.TypeOf(message), config.Topic)
		}
	}
}

// Push pushes messages to the underlying socket connection
// until the socket connection is closed.
func (c *Client) Push(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case command := <-c.commands:
			switch cmd := command.(type) {
			case send:
				config := cmd.message.GetConfig()
				if _, exists := c.codec.GetConfig(config.Kind); !exists {
					c.config.Log.Errorf("No MessageConfig associated with packet of kind %v\n", config.Kind)
					continue
				}

				packet := &Packet{Kind: config.Kind, Bytes: cmd.message.Marshal()}
				if err := JoinPacket(c.writer, packet); err != nil {
					// disconnected. when the writer reaches its capacity, it
					// automatically performs a flush to be able to write more
					// bytes, which may produce an error if the socket connection
					// is dropped
					return
				}

			case flush:
				if err := c.writer.Flush(); err != nil {
					return
				}

			case terminate:
				if err := c.connection.Close(); err != nil {
					return
				}
			}
		}
	}
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
