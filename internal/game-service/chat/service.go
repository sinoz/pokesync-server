package chat

import (
	"reflect"

	"gitlab.com/pokesync/game-service/internal/game-service/character"
	"gitlab.com/pokesync/game-service/internal/game-service/client"
	"go.uber.org/zap"
)

// Config holds configurations specifically for the chat service.
type Config struct {
	SessionConfig SessionConfig
	WorkerCount   int
}

// Service is an implementation of a public chat service and provides
// chatting capabilities for users across different channels.
type Service struct {
	config Config
	logger *zap.SugaredLogger

	routing *client.Router
	mailbox client.Mailbox
}

const (
	ServiceConnectTopic = "connect_to_chat_service"
	CreateChannelTopic  = "create_channel"
	JoinChannelTopic    = "join_channel"
	RemoveChannelTopic  = "remove_channel"
)

// ConnectToChatService is a request for a client to be connected
// to the chat service to start chatting.
type ConnectToChatService struct {
	DisplayName character.DisplayName
	UserGroup   character.UserGroup
}

// CreateChannel is a request to create and add a new channel
// for users to join.
type CreateChannel struct {
	Topic string
}

// JoinChannel is a request of a user to join an existing channel.
type JoinChannel struct {
	Topic    string
	ClientID client.ID
}

// RemoveChannel is a request to remove a channel from existence.
type RemoveChannel struct {
	Topic string
}

// messageTopicsOfInterest is a slice of message Topic's that the game
// Service has any interest in for processing.
var messageTopicsOfInterest = []client.Topic{
	ServiceConnectTopic,
	CreateChannelTopic,
	JoinChannelTopic,
	RemoveChannelTopic,
	SelectChatChannelConfig.Topic,
	SubmitChatMessageConfig.Topic,
	client.TerminationTopic,
}

// NewService constructs a new chat Service.
func NewService(config Config, logger *zap.SugaredLogger, routing *client.Router) *Service {
	service := &Service{
		config:  config,
		logger:  logger,
		routing: routing,
	}

	service.mailbox = routing.CreateMailbox()
	for _, topic := range messageTopicsOfInterest {
		routing.SubscribeMailboxToTopic(topic, service.mailbox)
	}

	go service.receive()

	return service
}

// receive receives and handles messages from the specified mailbox.
func (service *Service) receive() {
	for {
		select {
		case mail := <-service.mailbox:
			service.handleMail(mail)
		}
	}
}

// handleMail handles the given client Mail.
func (service *Service) handleMail(mail client.Mail) {
	switch message := mail.Payload.(type) {
	case ConnectToChatService:
		// TODO

	case CreateChannel:
		// TODO

	case JoinChannel:
		// TODO

	case RemoveChannel:
		// TODO

	case SelectChatChannel:
		// TODO

	case SubmitChatMessage:
		// TODO

	case client.Terminated:
		// TODO

	default:
		service.logger.Errorf("unexpected message received of type %v", reflect.TypeOf(message))
	}
}

// Stop stops this Service and cleans up resources.
func (service *Service) Stop() {
	// TODO stop all sessions

	close(service.mailbox)
}
