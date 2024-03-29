package game

import (
	"context"
	"reflect"
	"time"

	"gitlab.com/pokesync/game-service/pkg/event"

	"gitlab.com/pokesync/game-service/internal/game-service/chat"
	"gitlab.com/pokesync/game-service/internal/game-service/game/transport"

	"gitlab.com/pokesync/game-service/internal/game-service/game/entity"

	"gitlab.com/pokesync/game-service/internal/game-service/account"
	"gitlab.com/pokesync/game-service/internal/game-service/character"
	"gitlab.com/pokesync/game-service/internal/game-service/client"
	"go.uber.org/zap"
)

// Config holds configurations specific to the game service.
type Config struct {
	IntervalRate time.Duration

	EntityLimit int

	CharacterFetchTimeout time.Duration

	SessionConfig SessionConfig

	ClockRate         time.Duration
	ClockSynchronizer ClockSynchronizer

	Modules []Module
}

// CharacterProvider attempts to provide a character Profile.
type CharacterProvider func(email account.Email) <-chan character.LoadResult

// CharacterSaver attempts to save character Profile's.
type CharacterSaver func(email account.Email, profile *character.Profile)

// pulse represents a tick or a single heartbeat.
type pulse struct{}

// pulseInstance is a cached instance of the gamePulse type.
var pulseInstance = pulse{}

// pulser sends heartbeat-like pulses into the game to
// continuously process the game world.
type pulser struct {
	rate     time.Duration
	lastTime time.Time

	quit chan bool

	pulses chan pulse
	resume chan time.Duration
}

// Service is an implementation of a game service that provides
// gameplay capabilities to logged in users.
type Service struct {
	config Config
	assets *AssetBundle

	logger *zap.SugaredLogger

	routing  *client.Router
	sessions *SessionRegistry

	mailbox client.Mailbox
	pulser  *pulser

	characterProvider CharacterProvider
	characterSaver    CharacterSaver

	game *Game
}

// Game represents the game, mkay.
type Game struct {
	world         *entity.World
	entityFactory *EntityFactory
	eventBus      event.Bus
	grid          *Grid
	chatCommands  *ChatCommandRegistry
}

const (
	// AuthenticationEventTopic is a topic for authentication events.
	AuthenticationEventTopic = "auth_event"
)

// messageTopicsOfInterest is a slice of message Topic's that the game
// Service has any interest in for processing.
var messageTopicsOfInterest = []client.Topic{
	AuthenticationEventTopic,
	transport.AttachFollowerConfig.Topic,
	transport.SwitchPartySlotsConfig.Topic,
	transport.ChangeMovementTypeConfig.Topic,
	transport.MoveAvatarConfig.Topic,
	transport.ClearFollowerConfig.Topic,
	transport.ClickTeleportConfig.Topic,
	transport.ContinueDialogueConfig.Topic,
	transport.FaceDirectionConfig.Topic,
	transport.InteractWithEntityConfig.Topic,
	transport.SelectPlayerOptionConfig.Topic,
	transport.SubmitChatCommandConfig.Topic,
	client.TerminationTopic,
}

// Authenticated is an event of a user having been authenticated
// and is ready to have their character profile fetched so that
// they can be registered into the game.
type Authenticated struct {
	Account account.Account
}

// CharacterLoaded is an event of a client user having their character
// profile loaded and is thus ready to be registered into the game world.
type CharacterLoaded struct {
	Account   account.Account
	Character *character.Profile
}

// NewService constructs a new game Service.
func NewService(config Config, routing *client.Router, characterProvider CharacterProvider, characterSaver CharacterSaver, assets *AssetBundle, logger *zap.SugaredLogger) *Service {
	service := &Service{
		config: config,

		logger: logger,

		assets: assets,

		characterProvider: characterProvider,
		characterSaver:    characterSaver,

		routing: routing,
	}

	service.sessions = NewSessionRegistry()
	service.game = NewGame(config, assets, logger)

	service.pulser = newPulser(config.IntervalRate)
	service.mailbox = routing.CreateMailbox()

	for _, topic := range messageTopicsOfInterest {
		routing.SubscribeMailboxToTopic(topic, service.mailbox)
	}

	for _, module := range config.Modules {
		module(&DependencyKit{
			assets: assets,
			game:   service.game,
			logger: logger,
		})
	}

	go service.receive()

	return service
}

// NewGame constructs a new Game.
func NewGame(config Config, assets *AssetBundle, logger *zap.SugaredLogger) *Game {
	world := entity.NewWorld(config.EntityLimit)
	entityFactory := NewEntityFactory(world, assets)
	eventBus := event.NewSerialBus()
	chatCommands := NewChatCommandRegistry()

	game := &Game{
		world:         world,
		entityFactory: entityFactory,
		eventBus:      eventBus,
		grid:          assets.Grid,
		chatCommands:  chatCommands,
	}

	world.AddSystem(NewInboundNetworkSystem(
		logger,

		withAttachFollowerHandler(attachFollower()),
		withClearFollowerHandler(clearFollower()),
		withSwitchPartySlotHandler(switchPartySlots()),
		withClickTeleportHandler(clickTeleport()),
		withContinueDialogueHandler(continueDialogue()),
		withSelectPlayerOptionHandler(selectPlayerOption()),
		withEntityInteraction(interactWithEntity()),
		withDirectionFacingHandler(faceDirection()),
		withMoveAvatarHandler(moveAvatar()),
		withMovementTypeChangeHandler(changeMovementType()),
		withSubmitChatCommandHandler(submitChatCommand(chatCommands)),
	))

	world.AddSystem(NewWalkingSystem(game.grid, AStarRouteFinder()))
	world.AddSystem(NewRunningSystem(game.grid, AStarRouteFinder()))
	world.AddSystem(NewCyclingSystem(game.grid, AStarRouteFinder()))
	world.AddSystem(NewDayNightSystem(config.ClockRate, config.ClockSynchronizer))
	world.AddSystem(NewMapViewSystem(game.grid))
	world.AddSystem(NewOutboundNetworkSystem())

	return game
}

// newPulser constructs a new pulser that operates at the specified
// interval rate, specified in milliseconds.
func newPulser(intervalRate time.Duration) *pulser {
	pulser := new(pulser)

	pulser.rate = intervalRate
	pulser.lastTime = time.Now()

	pulser.quit = make(chan bool, 1)
	pulser.resume = make(chan time.Duration)
	pulser.pulses = make(chan pulse)

	go func() {
		for {
			if pulser.hasQuitPulsing() {
				break
			}

			// send out a heartbeat into the stream
			pulser.pulses <- pulseInstance

			// and then we wait for confirmation to resume, which is to
			// indicate that processing of the pulse has been completed.
			// the 'resume' channel gives us back the amount of time it
			// had taken to process the pulse, which we can then use to
			// calculate the amount of time we have left to sleep until
			// we send out the next heartbeat.
			timeElapsed := <-pulser.resume
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
		go service.onAuthenticated(mail.Context, mail.Client, message.Account)

	case CharacterLoaded:
		service.onCharacterLoaded(mail.Client, message.Account, message.Character)

	case client.Message:
		session := service.sessions.Get(mail.Client.ID)
		if session == nil {
			return
		}

		session.QueueCommand(message)

	case client.Terminated:
		session := service.sessions.Remove(mail.Client.ID)
		if session == nil {
			return
		}

		service.game.RemovePlayer(session.Player)
		service.characterSaver(session.Email, service.transformPlayerToCharacterProfile(session.Player))

	default:
		service.logger.Errorf("unexpected message received of type %v", reflect.TypeOf(message))
	}
}

// transformPlayerToCharacterProfile transforms the given Player instance
// into a character Profile that can then be persisted.
func (service *Service) transformPlayerToCharacterProfile(player *Player) *character.Profile {
	displayName := player.DisplayName()
	userGroup := player.Rank()
	position := player.Position()
	bicycleType := player.BicycleType()
	coinBag := player.CoinBag()

	lastLoggedIn := time.Now()

	return &character.Profile{
		DisplayName:  displayName,
		LastLoggedIn: &lastLoggedIn,
		UserGroup:    userGroup,

		Gender:      0,
		BicycleType: int(bicycleType),

		PokeDollars:   coinBag.Dollars,
		DonatorPoints: coinBag.DonatorPoints,

		MapX:   position.MapX,
		MapZ:   position.MapZ,
		LocalX: position.LocalX,
		LocalZ: position.LocalZ,
	}
}

// onAuthenticated reacts to the given Client user having been authenticated.
func (service *Service) onAuthenticated(ctx context.Context, cl *client.Client, account account.Account) {
	select {
	case <-ctx.Done():
		return

	case result := <-service.characterProvider(account.Email):
		if result.Error != nil {
			service.logger.Error(result.Error)

			cl.SendNow(&transport.UnableToFetchProfile{})
			cl.Terminate()

			return
		}

		if result.Profile == nil {
			cl.SendNow(&transport.UnableToFetchProfile{})
			cl.Terminate()

			return
		}

		characterLoadEvent := CharacterLoaded{Account: account, Character: result.Profile}
		mail := client.Mail{Context: ctx, Client: cl, Payload: characterLoadEvent}

		service.mailbox <- mail

	case <-time.After(service.config.CharacterFetchTimeout):
		cl.SendNow(&transport.RequestTimedOut{})
		cl.Terminate()
	}
}

// onCharacterLoaded reacts to the given Client user having a character
// profile loaded for him/her.
func (service *Service) onCharacterLoaded(cl *client.Client, account account.Account, character *character.Profile) {
	position := Position{
		MapX:   character.MapX,
		MapZ:   character.MapZ,
		LocalX: character.LocalX,
		LocalZ: character.LocalZ,
	}

	plr := service.game.CreatePlayer(position, Gender(character.Gender), character.DisplayName, character.UserGroup)
	if ok := service.game.AddPlayer(plr); !ok {
		cl.SendNow(&transport.WorldFull{})
		cl.Terminate()

		return
	}

	session := NewInstalledSession(cl, service.config.SessionConfig, account.Email, plr)
	service.sessions.Put(cl.ID, session)

	plr.Add(&SessionComponent{session: session})

	plr.CoinBag().AddPokeDollars(character.PokeDollars)
	plr.CoinBag().AddDonatorPoints(character.DonatorPoints)

	plr.SetBicycleType(BicycleType(character.BicycleType))

	cl.SendNow(&transport.LoginSuccess{
		PID:         uint16(plr.ID),
		DisplayName: string(character.DisplayName),

		Gender:    byte(character.Gender),
		UserGroup: byte(character.UserGroup),

		MapX:   uint16(character.MapX),
		MapZ:   uint16(character.MapZ),
		LocalX: uint16(character.LocalX),
		LocalZ: uint16(character.LocalZ),
	})

	service.routing.Publish(chat.ServiceConnectTopic, client.Mail{
		Client: cl,
		Payload: chat.ConnectToChatService{
			DisplayName: character.DisplayName,
		},
	})
}

// CreatePlayer creates a new Player-like Entity.
func (game *Game) CreatePlayer(position Position, gender Gender, displayName character.DisplayName, userGroup character.UserGroup) *Player {
	return PlayerBy(game.entityFactory.CreatePlayer(position, gender, displayName, userGroup))
}

// CreateNpc creates a new Npc-like Entity.
func (game *Game) CreateNpc(modelID ModelID, position Position) *Npc {
	return NpcBy(game.entityFactory.CreateNpc(position, modelID))
}

// CreateMonster creates a new Monster-like Entity.
func (game *Game) CreateMonster(position Position, data MonsterData) *Monster {
	return MonsterBy(game.entityFactory.CreateMonster(position, data.ModelID), data)
}

// AddPlayer attempts to add the given Player to the game world.
func (game *Game) AddPlayer(plr *Player) bool {
	return game.AddEntity(plr.Entity)
}

// RemovePlayer removes the given Player-like entity.
func (game *Game) RemovePlayer(plr *Player) {
	game.RemoveEntity(plr.Entity)
}

// AddNpc attempts to add the given Npc to the game world.
func (game *Game) AddNpc(npc *Npc) bool {
	return game.AddEntity(npc.Entity)
}

// RemoveNpc removes the given Npc-like entity.
func (game *Game) RemoveNpc(npc *Npc) {
	game.RemoveEntity(npc.Entity)
}

// AddMonster attempts to add the given Monster to the game world.
func (game *Game) AddMonster(monster *Monster) bool {
	return game.AddEntity(monster.Entity)
}

// RemoveMonster removes the given Monster-like entity.
func (game *Game) RemoveMonster(monster *Monster) {
	game.RemoveEntity(monster.Entity)
}

// AddEntity attempts to add the given Entity to the game world.
func (game *Game) AddEntity(entity *entity.Entity) bool {
	return game.world.AddEntity(entity)
}

// RemoveEntity removes the given Entity from the game world.
func (game *Game) RemoveEntity(entity *entity.Entity) {
	game.world.DestroyEntity(entity)
}

// pulse is called every game pulse to process the game.
func (game *Game) pulse(deltaTime time.Duration) error {
	return game.world.Update(deltaTime)
}

// pulse is called every pulse or tick to process the game.
func (service *Service) pulse() {
	deltaTime := time.Since(service.pulser.lastTime)
	service.pulser.lastTime = time.Now()

	if err := service.game.pulse(deltaTime); err != nil {
		service.logger.Error(err)
	}

	service.pulser.resume <- time.Since(service.pulser.lastTime)
}

// Stop stops this Service and cleans up resources.
func (service *Service) Stop() {
	service.pulser.quitPulsing()
}
