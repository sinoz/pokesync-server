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

// PulseTask is the task to execute every pulse.
type pulseTask func(time.Duration)

// pulser sends heartbeat-like pulses into the game to
// continuously process the game world.
type pulser struct {
	isRunning bool
	rate      time.Duration
	deltaTime time.Duration
	lastTime  time.Time
}

// Service is an implementation of a game service that provides
// gameplay capabilities to logged in users.
type Service struct {
	config Config

	assets *AssetBundle

	characters character.Repository

	routing *client.Router

	sessions *session.Registry

	pulser *pulser
	game   *Game
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

	service.pulser = newPulser(config.IntervalRate, service.pulse)

	mailbox := routing.Subscribe(AuthenticationEventTopic)
	go service.receiver(mailbox)

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
func newPulser(intervalRate time.Duration, runTask pulseTask) *pulser {
	pulser := new(pulser)

	pulser.rate = intervalRate
	pulser.isRunning = true
	pulser.lastTime = time.Now()

	go func() {
		for pulser.isRunning {
			pulser.deltaTime = time.Since(pulser.lastTime)
			pulser.lastTime = time.Now()

			runTask(pulser.deltaTime)

			timeElapsed := time.Since(pulser.lastTime)
			timeToSleep := intervalRate - timeElapsed
			if timeToSleep > 0 {
				time.Sleep(timeToSleep)
			}
		}
	}()

	return pulser
}

// receiver receives and handles client messages from the specified mailbox.
func (service *Service) receiver(mailbox client.Mailbox) {
	for mail := range mailbox {
		switch message := mail.Payload.(type) {
		case Authenticated:
			// TODO have a worker do this

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

		case client.Message:
			session := service.sessions.Get(mail.Client.ID)
			if session == nil {
				continue
			}

			session.QueueCommand(message)

		default:
			service.config.Logger.Errorf("unexpected message received of type %v", reflect.TypeOf(message))
		}
	}
}

// pulse is called every pulse or tick to process the game. The given
// delta time parameter is the amount of time that has elapsed since
// the last pulse.
func (service *Service) pulse(deltaTime time.Duration) {
	if err := service.game.pulse(deltaTime); err != nil {
		service.config.Logger.Error(err)
	}
}

// TearDown terminates this Service and cleans up resources.
func (service *Service) TearDown() {
	// TODO
}
