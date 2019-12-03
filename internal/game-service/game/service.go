package game

import (
	"reflect"
	"time"

	"gitlab.com/pokesync/game-service/internal/game-service/game/entity"
	"gitlab.com/pokesync/game-service/internal/game-service/game/session"

	"gitlab.com/pokesync/game-service/internal/game-service/account"
	"gitlab.com/pokesync/game-service/internal/game-service/character"
	"gitlab.com/pokesync/game-service/internal/game-service/client"
	"go.uber.org/zap"
)

// Config holds configurations specific to the game service.
type Config struct {
	IntervalRate time.Duration

	EntityLimit int

	SessionConfig session.Config

	Logger *zap.SugaredLogger

	ClockRate         time.Duration
	ClockSynchronizer ClockSynchronizer
}

// pulse represents a tick or a single heartbeat.
type pulse struct{}

// pulseInstance is a cached instance of the gamePulse type.
var pulseInstance = pulse{}

// pulser sends heartbeat-like pulses into the game to
// continuously process the game world.
type pulser struct {
	rate     time.Duration
	lastTime time.Time

	quit   chan bool
	pulses chan pulse
}

// Service is an implementation of a game service that provides
// gameplay capabilities to logged in users.
type Service struct {
	config Config
	assets *AssetBundle

	routing  *client.Router
	sessions *session.Registry

	mailbox client.Mailbox
	pulser  *pulser

	characters character.Repository

	game *Game
}

const (
	// AuthenticationEventTopic is a topic for authentication events.
	AuthenticationEventTopic = "auth_event"
)

// Authenticated is an event of a user having been authenticated
// and is ready to have their character profile fetched so that
// they can be registered into the game.
type Authenticated struct {
	Account account.Account
}

// NewService constructs a new game Service.
func NewService(config Config, routing *client.Router, characters character.Repository, assets *AssetBundle) *Service {
	service := &Service{
		config: config,

		assets: assets,

		characters: characters,

		routing: routing,
	}

	service.sessions = session.NewRegistry()
	service.game = NewGame(assets, config.EntityLimit)

	service.pulser = newPulser(config.IntervalRate)
	service.mailbox = routing.Subscribe(AuthenticationEventTopic)

	go service.receive()

	return service
}

// createWorld constructs a new instance of a World, preconfigured
// with all of its necessary system and processors for the game service
// to process game logic.
func createWorld(config Config, assets *AssetBundle) *entity.World {
	world := entity.NewWorld(config.EntityLimit)

	world.AddSystem(NewInboundNetworkSystem())
	world.AddSystem(NewDayNightSystem(config.ClockRate, config.ClockSynchronizer))
	world.AddSystem(NewMapViewSystem())
	world.AddSystem(NewOutboundNetworkSystem())

	return world
}

// newPulser constructs a new pulser that operates at the specified
// interval rate, specified in milliseconds.
func newPulser(intervalRate time.Duration) *pulser {
	pulser := new(pulser)

	pulser.rate = intervalRate
	pulser.lastTime = time.Now()

	pulser.quit = make(chan bool, 1)
	pulser.pulses = make(chan pulse)

	go func() {
		for {
			if pulser.hasQuitPulsing() {
				break
			}

			pulser.pulses <- pulseInstance

			timeElapsed := time.Since(pulser.lastTime)
			timeToSleep := intervalRate - timeElapsed
			if timeToSleep > 0 {
				time.Sleep(timeToSleep)
			}
		}
	}()

	return pulser
}

// quitPulsing sends a signal to have the pulser stop running.
func (pulser *pulser) quitPulsing() {
	pulser.quit <- true
	close(pulser.pulses)
}

// hasQuitPulsing returns whether a signal was fired to quit pulsing.
func (pulser *pulser) hasQuitPulsing() bool {
	select {
	case <-pulser.quit:
		return true
	default:
		return false
	}
}

// receive receives and handles messages from the specified mailbox
// and also deals with game pulses.
func (service *Service) receive() {
	for {
		select {
		case <-service.pulser.pulses:
			service.pulse()

		case mail := <-service.mailbox:
			service.handleMail(mail)
		}
	}
}

// handleMail handles the given client Mail.
func (service *Service) handleMail(mail client.Mail) {
	switch message := mail.Payload.(type) {
	case Authenticated:
		service.onAuthenticated(mail.Client, message.Account)

	case client.Message:
		session := service.sessions.Get(mail.Client.ID)
		if session == nil {
			return
		}

		session.QueueCommand(message)

	default:
		service.config.Logger.Errorf("unexpected message received of type %v", reflect.TypeOf(message))
	}
}

// onAuthenticated reacts to the given Client user having been authenticated.
func (service *Service) onAuthenticated(cl *client.Client, account account.Account) {
	// TODO

	// mail.Client.SendNow(&transport.LoginSuccess{
	// 	PID:         1,
	// 	DisplayName: string(profile.DisplayName),
	// 	Gender:      byte(profile.Gender),
	// 	UserGroup:   byte(profile.UserGroup),

	// 	MapX:   uint16(profile.MapX),
	// 	MapZ:   uint16(profile.MapZ),
	// 	LocalX: uint16(profile.LocalX),
	// 	LocalZ: uint16(profile.LocalZ),
	// })
}

// pulse is called every pulse or tick to process the game.
func (service *Service) pulse() {
	deltaTime := time.Since(service.pulser.lastTime)
	service.pulser.lastTime = time.Now()

	if err := service.game.pulse(deltaTime); err != nil {
		service.config.Logger.Error(err)
	}
}

// TearDown terminates this Service and cleans up resources.
func (service *Service) TearDown() {
	service.pulser.quitPulsing()
}
