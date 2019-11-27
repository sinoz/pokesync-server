package game

import (
	"fmt"
	"time"

	ecs "gitlab.com/pokesync/ecs/src"
	"gitlab.com/pokesync/game-service/internal/game-service/account"
	"gitlab.com/pokesync/game-service/internal/game-service/character"
	"gitlab.com/pokesync/game-service/internal/game-service/client"
)

// Unbounded is for parameters such as the job limit.
const Unbounded = -1

// Config holds configurations specific to the game service.
type Config struct {
	IntervalRate time.Duration
	JobLimit     int
	WorkerCount  int
	EntityLimit  int
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
	config     Config
	routing    *client.Router
	jobQueue   chan Job
	assets     *AssetBundle
	pulser     *pulser
	characters character.Repository
	world      *ecs.World
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

// Job is an interface for game specific jobs.
type Job interface{}

// NewService constructs a new game Service.
func NewService(config Config, routing *client.Router, characters character.Repository, assets *AssetBundle) *Service {
	var jobQueue chan Job
	if config.JobLimit == Unbounded {
		jobQueue = make(chan Job)
	} else {
		jobQueue = make(chan Job, config.JobLimit)
	}

	service := &Service{
		config:     config,
		routing:    routing,
		jobQueue:   jobQueue,
		assets:     assets,
		characters: characters,
	}

	service.world = createWorld(config.EntityLimit, assets)
	service.pulser = newPulser(config.IntervalRate, service.pulse)

	mailbox := routing.Subscribe(AuthenticationEventTopic)
	service.receiver(mailbox)

	for i := 0; i < config.WorkerCount; i++ {
		service.spawnWorker()
	}

	return service
}

// createWorld constructs a new instance of a World, preconfigured
// with all of its necessary system and processors for the game service
// to process game logic.
func createWorld(entityLimit int, assets *AssetBundle) *ecs.World {
	world := ecs.NewWorld(entityLimit)

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
			timeToSleep := (intervalRate * time.Millisecond) - timeElapsed
			if timeToSleep > 0 {
				time.Sleep(timeToSleep)
			}
		}
	}()

	return pulser
}

// receiver receives and handles client messages from the specified mailbox.
func (service *Service) receiver(mailbox client.Mailbox) {
	go func() {
		for mail := range mailbox {
			switch message := mail.Payload.(type) {
			case Authenticated:
				// TODO have a worker do this
				profile, err := service.characters.Get(message.Account.Email)
				if err != nil {
					mail.Client.SendNow(&UnableToFetchProfile{})
					mail.Client.Terminate()

					continue
				}

				mail.Client.SendNow(&LoginSuccess{
					PID:         1,
					DisplayName: string(profile.DisplayName),
					Gender:      byte(profile.Gender),
					UserGroup:   byte(profile.UserGroup),

					MapX:   uint16(profile.MapX),
					MapZ:   uint16(profile.MapZ),
					LocalX: uint16(profile.LocalX),
					LocalZ: uint16(profile.LocalZ),
				})
			}
		}
	}()
}

// spawnWorker spawns a worker goroutine that reads from the
// service's job queue.
func (service *Service) spawnWorker() {
	go func() {
		for job := range service.jobQueue {
			fmt.Println(job)
		}
	}()
}

// pulse is called every pulse or tick to process the game. The given
// delta time parameter is the amount of time that has elapsed since
// the last pulse.
func (service *Service) pulse(deltaTime time.Duration) {
	service.world.Update(deltaTime)
}
