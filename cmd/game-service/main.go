package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"gitlab.com/pokesync/game-service/internal/game-service/account"
	"gitlab.com/pokesync/game-service/internal/game-service/character"
	"gitlab.com/pokesync/game-service/internal/game-service/chat"
	"gitlab.com/pokesync/game-service/internal/game-service/client"
	"gitlab.com/pokesync/game-service/internal/game-service/discord"
	"gitlab.com/pokesync/game-service/internal/game-service/game"
	gameTransport "gitlab.com/pokesync/game-service/internal/game-service/game/transport"
	"gitlab.com/pokesync/game-service/internal/game-service/login"
	"gitlab.com/pokesync/game-service/internal/game-service/server"
	"gitlab.com/pokesync/game-service/internal/game-service/status"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ClientBuildNo is the build number of the client this game
// server is supporting.
const ClientBuildNo = client.BuildNumber(1)

// WorldIDEnv is the name of the environment variable of the world id.
const WorldIDEnv = "POKESYNC_WORLD_ID"

// TCPServerHostEnv is the name of the environment variable of
// the server host.
const TCPServerHostEnv = "POKESYNC_HOST"

// TCPServerPortEnv is the name of the environment variable of
// the server port.
const TCPServerPortEnv = "POKESYNC_PORT"

// RedisHostEnv is the name of the environment variable of
// the host of the Redis server to connect to.
const RedisHostEnv = "POKESYNC_REDIS_HOST"

// RedisPortEnv is the name of the environment variable of
// the port of the Redis server to connect to.
const RedisPortEnv = "POKESYNC_REDIS_PORT"

// DefaultWorldID is the id of the game world to fallback to if no environment
// variable is set.
const DefaultWorldID = 1

// DefaultTCPServerHost is the host to fallback to if no environment
// variable is set.
const DefaultTCPServerHost = "localhost"

// DefaultTCPServerPort is the port to fallback to if no environment
// variable is set.
const DefaultTCPServerPort = 23192

// DefaultRedisHost is the host to fallback to if no environment
// variable is set.
const DefaultRedisHost = "localhost"

// DefaultRedisPort is the port to fallback to if no environment
// variable is set.
const DefaultRedisPort = 6379

// loginCodec is a message Codec that holds marshallers and demarshallers
// specific for the login aspect of the server.
var loginCodec = client.NewCodec().
	Include(login.RequestConfig).
	Include(login.RequestTimedOutConfig).
	Include(login.ErrorDuringAccountFetchConfig).
	Include(login.AccountDisabledConfig).
	Include(login.AlreadyLoggedInConfig).
	Include(login.InvalidCredentialsConfig)

// chatCodec is a message Codec that holds marshallers and demarshallers
// specific for the public chatting aspect of the server.
var chatCodec = client.NewCodec().
	Include(chat.DisplayChatMessageConfig).
	Include(chat.SelectChatChannelConfig).
	Include(chat.SubmitChatMessageConfig).
	Include(chat.SwitchChatChannelConfig)

// gameCodec is a message Codec that holds marshallers and demarshallers
// specific for the game aspect of the server.
var gameCodec = client.NewCodec().
	Include(gameTransport.UnableToFetchProfileConfig).
	Include(gameTransport.LoginSuccessConfig).
	Include(gameTransport.RefreshMapConfig).
	Include(gameTransport.MoveAvatarConfig).
	Include(gameTransport.MoveCameraConfig).
	Include(gameTransport.ResetCameraConfig).
	Include(gameTransport.AttachFollowerConfig).
	Include(gameTransport.ClearFollowerConfig).
	Include(gameTransport.ChangeMovementTypeConfig).
	Include(gameTransport.ClickTeleportConfig).
	Include(gameTransport.CloseDialogueConfig).
	Include(gameTransport.ContinueDialogueConfig).
	Include(gameTransport.FaceDirectionConfig).
	Include(gameTransport.EntityUpdateConfig).
	Include(gameTransport.InteractWithEntityConfig).
	Include(gameTransport.SubmitChatCommandConfig).
	Include(gameTransport.SelectCharacterConfig).
	Include(gameTransport.SetDonatorPointsConfig).
	Include(gameTransport.SetPokeDollarsConfig).
	Include(gameTransport.SetPartySlotConfig).
	Include(gameTransport.SwitchPartySlotsConfig).
	Include(gameTransport.SelectPlayerOptionConfig).
	Include(gameTransport.SetServerTimeConfig)

// messageCodec holds demarshallers and marshallers of messages.
var messageCodec = client.NewCodec().
	Join(loginCodec).
	Join(chatCodec).
	Join(gameCodec)

// The main entry point to this game server application.
func main() {
	logger, err := createZapLogger()
	if err != nil {
		log.Fatal("Failed to create zap logger", err)
	}

	logger.Info("Starting PokeSync ...")

	worldID := getWorldIDFromEnv()

	tcpHost := getTCPServerHostFromEnv()
	tcpPort := getTCPServerPortFromEnv()

	redisHost := getRedisHostFromEnv()
	redisPort := getRedisPortFromEnv()

	redisClient, err := connectToRedis(redisHost, redisPort)
	if err != nil {
		log.Fatal("Failed to connect to a Redis server instance", err)
	}

	logger.Infof("Connected to Redis instance at %v:%v", redisHost, redisPort)

	assetsConfig := game.AssetConfig{
		ItemDirectory:    "assets/config/item",
		NpcDirectory:     "assets/config/npc",
		MonsterDirectory: "assets/config/monster",
		ObjectDirectory:  "assets/config/object",
		WorldDirectory:   "assets/config/world",
	}

	assetBundle, err := game.LoadAssetBundle(assetsConfig)
	if err != nil {
		logger.Fatal(err)
	}

	routingConfig := client.RouterConfig{
		PublicationTimeout: 1 * time.Second,
	}

	routing := client.NewRouter(routingConfig)

	accountConfig := account.Config{
		WorkerCount: runtime.NumCPU(),
	}

	characterCacheConfig := character.CacheConfig{
		ExpireAfter: 5 * time.Minute,
	}

	charactersConfig := character.Config{
		WorkerCount: runtime.NumCPU(),
	}

	gameConfig := game.Config{
		IntervalRate:          50 * time.Millisecond,
		CharacterFetchTimeout: 5 * time.Second,
		EntityLimit:           32768,

		// not to be confused with IntervalRate, ClockRate is the rate at
		// which the game time progresses (think day/night, seasons) and
		// is not actually the tick rate. The ClockRate can also be seen
		// as the ratio of game time to real time, which is currently 4:1
		ClockRate:         250 * time.Millisecond,
		ClockSynchronizer: game.NewGMT0Synchronizer(),

		SessionConfig: game.SessionConfig{
			CommandLimit: 16,
			EventLimit:   256,
		},
	}

	statusConfig := status.Config{
		RefreshRate: 15 * time.Second,
	}

	discordConfig := discord.Config{}

	authConfig := login.AuthConfig{
		AccountFetchTimeout: 5 * time.Second,
	}

	loginConfig := login.Config{
		WorkerCount: runtime.NumCPU(),
	}

	chatConfig := chat.Config{
		WorkerCount: runtime.NumCPU(),

		SessionConfig: chat.SessionConfig{
			BufferLimit: 32,
		},
	}

	clientConfig := client.Config{
		MessageCodec:    *messageCodec,
		ReadBufferSize:  512,
		WriteBufferSize: 2048,
		CommandLimit:    32,
	}

	serverConfig := server.Config{
		ClientConfig: clientConfig,
	}

	accountRepository := account.NewInMemoryRepository()
	passwordMatcher := account.BasicPasswordMatcher()

	characterCache := character.NewRedisCache(characterCacheConfig, redisClient)
	characterRepository := character.NewInMemoryRepository()
	characterService := character.NewService(charactersConfig, logger, characterCache, characterRepository)

	accountService := account.NewService(accountConfig, logger, accountRepository)
	chatService := chat.NewService(chatConfig, logger, routing)

	authenticator := login.NewAuthenticator(
		authConfig,
		accountService.LoadAccount,
		passwordMatcher,
	)

	loginService := login.NewService(loginConfig, logger, authenticator, routing)

	gameService := game.NewService(gameConfig, routing, characterService.LoadProfile, characterService.SaveProfile, assetBundle, logger)
	discordService := discord.NewService(discordConfig, logger)
	statusService := status.NewService(statusConfig, logger, status.NewRedisNotifier(redisClient, worldID), status.NewProvider(gameService))

	// should something go wrong and cause a panic, always safely
	// tear down these services
	defer func() {
		accountService.Stop()
		chatService.Stop()
		loginService.Stop()
		discordService.Stop()
		gameService.Stop()
		statusService.Stop()
	}()

	logger.Info("Client build: ", ClientBuildNo)
	logger.Info("World ID: ", worldID)

	logger.Info("Item configs loaded: ", assetBundle.Items.Count())
	logger.Info("Npc configs loaded: ", assetBundle.Npcs.Count())
	logger.Info("Monster configs loaded: ", assetBundle.Monsters.Count())

	logger.Info("Game pulse rate: ", gameConfig.IntervalRate)
	logger.Info("Game clock rate: ", gameConfig.ClockRate)

	logger.Info("World entity limit: ", gameConfig.EntityLimit)

	logger.Info("Account worker count: ", accountConfig.WorkerCount)
	logger.Info("Login worker count: ", loginConfig.WorkerCount)
	logger.Info("Character worker count: ", charactersConfig.WorkerCount)
	logger.Info("Chat worker count: ", chatConfig.WorkerCount)

	logger.Info("Account fetch timeout: ", authConfig.AccountFetchTimeout)
	logger.Info("Character fetch timeout: ", gameConfig.CharacterFetchTimeout)

	logger.Info("Upstream byte limit: ", clientConfig.ReadBufferSize)
	logger.Info("Downstream byte limit: ", clientConfig.WriteBufferSize)

	logger.Info("Chat message buffer limit: ", chatConfig.SessionConfig.BufferLimit)

	logger.Info("Game Session command limit: ", gameConfig.SessionConfig.CommandLimit)
	logger.Info("Game Session event limit: ", gameConfig.SessionConfig.EventLimit)

	logger.Info("Server status update rate: ", statusConfig.RefreshRate)

	logger.Info("Router publication timeout: ", routingConfig.PublicationTimeout)

	tcpListener := server.NewTCPListener(serverConfig, routing, logger)
	if err := tcpListener.Bind(tcpHost, tcpPort); err != nil {
		logger.Fatal(err)
	}
}

func connectToRedis(host string, port int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprint(host, ":", port),
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func createZapLogger() (*zap.SugaredLogger, error) {
	conf := zap.NewDevelopmentConfig()

	conf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	conf.DisableStacktrace = true
	conf.DisableCaller = true

	logger, err := conf.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

func getWorldIDFromEnv() int {
	port, err := strconv.Atoi(os.Getenv(WorldIDEnv))
	if err != nil {
		return DefaultWorldID
	}

	return port
}

func getRedisHostFromEnv() string {
	host := os.Getenv(RedisHostEnv)
	if len(host) == 0 {
		return DefaultRedisHost
	}

	return host
}

func getRedisPortFromEnv() int {
	port, err := strconv.Atoi(os.Getenv(RedisPortEnv))
	if err != nil {
		return DefaultRedisPort
	}

	return port
}

func getTCPServerHostFromEnv() string {
	host := os.Getenv(TCPServerHostEnv)
	if len(host) == 0 {
		return DefaultTCPServerHost
	}

	return host
}

func getTCPServerPortFromEnv() int {
	port, err := strconv.Atoi(os.Getenv(TCPServerPortEnv))
	if err != nil {
		return DefaultTCPServerPort
	}

	return port
}
