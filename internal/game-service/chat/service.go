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

	sessions *SessionRegistry
}

const (
	ServiceConnectTopic = "connect_to_chat_service"
)

// ConnectToChatService is a request for a client to be connected
// to the chat service to start chatting.
type ConnectToChatService struct {
	DisplayName character.DisplayName
}

// messageTopicsOfInterest is a slice of message Topic's that the game
// Service has any interest in for processing.
var messageTopicsOfInterest = []client.Topic{
	ServiceConnectTopic,
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

	service.sessions = NewSessionRegistry()
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
		session := service.sessions.Get(mail.Client.ID)
		if session != nil {
			return
		}

		session = NewSession(mail.Client, service.config.SessionConfig)
		service.sessions.Put(mail.Client.ID, session)

	case client.Message:
		session := service.sessions.Get(mail.Client.ID)
		if session == nil {
			return
		}

		session.Enqueue(message)

	case client.Terminated:
		session := service.sessions.Remove(message.ID)
		if session == nil {
			service.logger.Warnf("attempted to remove a chat session from its registry by id %v but it didn't exist", message.ID)
			return
		}

		session.Stop()

	default:
		service.logger.Errorf("unexpected message received of type %v", reflect.TypeOf(message))
	}
}

// Stop stops this Service and cleans up resources.
func (service *Service) Stop() {
	// TODO stop all sessions

	close(service.mailbox)
}
