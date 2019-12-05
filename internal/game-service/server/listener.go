package server

import (
	"context"
	"net"
	"strconv"

	"gitlab.com/pokesync/game-service/internal/game-service/client"
	"go.uber.org/zap"
)

// Config holds configurations specific to the server listener.
type Config struct {
	ClientConfig client.Config
	Logger       *zap.SugaredLogger
}

// Listener listens for incoming connections at a port.
type Listener interface {
	Bind(address string, port int) error
	Unbind() error
}

// TcpListener is a Listener that listens at a TCP port.
type TcpListener struct {
	config   Config
	listener net.Listener
	Router   *client.Router
}

// NewTcpListener constructs a TcpListener.
func NewTcpListener(config Config, routing *client.Router) *TcpListener {
	return &TcpListener{
		config: config,
		Router: routing,
	}
}

// Bind binds the TcpListener to listen at the specified address and port.
func (l *TcpListener) Bind(address string, port int) error {
	fullAddress := address + ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", fullAddress)
	if err != nil {
		return err
	}

	l.listener = listener
	l.config.Logger.Info("Channel bound at: ", fullAddress)

	background := context.Background()

	for {
		connection, err := listener.Accept()
		if err != nil {
			return err
		}

		c := client.NewClient(connection, l.config.ClientConfig, l.Router)
		ctx, cancel := context.WithCancel(background)

		go c.Pull(ctx, cancel)
		go c.Push(ctx)
	}
}

// Unbind unbinds the TcpListener from listening at a port.
func (l *TcpListener) Unbind() error {
	return l.listener.Close()
}
