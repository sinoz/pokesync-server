package server

import (
	"context"
	"io"
	"net"
	"strconv"

	"gitlab.com/pokesync/game-service/internal/game-service/client"
	"go.uber.org/zap"
)

// Config holds configurations specific to the server listener.
type Config struct {
	ClientConfig client.Config
}

// Listener listens for incoming connections at a port.
type Listener interface {
	Bind(address string, port int) error
	Unbind() error
}

// TCPListener is a Listener that listens at a TCP port.
type TCPListener struct {
	config   Config
	listener net.Listener
	Router   *client.Router
	logger   *zap.SugaredLogger
}

// NewTCPListener constructs a TcpListener.
func NewTCPListener(config Config, routing *client.Router, logger *zap.SugaredLogger) *TCPListener {
	return &TCPListener{
		config: config,
		Router: routing,
		logger: logger,
	}
}

// Bind binds the TCPListener to listen at the specified address and port.
func (l *TCPListener) Bind(address string, port int) error {
	fullAddress := address + ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", fullAddress)
	if err != nil {
		return err
	}

	l.listener = listener
	l.logger.Info("Channel bound at: ", fullAddress)

	background := context.Background()

	for {
		connection, err := listener.Accept()
		if err != nil {
			return err
		}

		c := client.NewClient(connection, l.config.ClientConfig)
		ctx, cancel := context.WithCancel(background)

		go func() {
			for {
				if err := c.Pull(ctx, cancel, l.Router); err != nil {
					if err != io.EOF {
						l.logger.Error(err)
					}

					break
				}
			}
		}()

		go func() {
			for {
				if err := c.Push(ctx); err != nil {
					l.logger.Error(err)
					break
				}
			}
		}()
	}
}

// Unbind unbinds the TcpListener from listening at a port.
func (l *TCPListener) Unbind() error {
	return l.listener.Close()
}
