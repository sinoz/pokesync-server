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
	"gitlab.com/pokesync/game-service/internal/game-service/login"
	"gitlab.com/pokesync/game-service/internal/game-service/server"
	"gitlab.com/pokesync/game-service/internal/game-service/status"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ClientBuildNo is the build number of the client this game
// server is supporting.
const ClientBuildNo = client.BuildNumber(1)

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

// gameCodec is a message Codec that holds marshallers and demarshallers
// specific for the game aspect of the server.
var gameCodec = client.NewCodec().
	Include(game.UnableToFetchProfileConfig).
	Include(game.LoginSuccessConfig).
	Include(game.RefreshMapConfig).
	Include(game.MoveAvatarConfig).
	Include(game.MoveCameraConfig).
	Include(game.ResetCameraConfig).
	Include(game.AttachFollowerConfig).
	Include(game.ClearFollowerConfig).
	Include(game.ChangeMovementTypeConfig).
	Include(game.ClickTeleportConfig).
	Include(game.CloseDialogueConfig).
	Include(game.ContinueDialogueConfig).
	Include(game.DisplayChatMessageConfig).
	Include(game.FaceDirectionConfig).
	Include(game.EntityUpdateConfig).
	Include(game.InteractWithEntityConfig).
	Include(game.SelectCharacterConfig).
	Include(game.SelectChatChannelConfig).
	Include(game.SetDonatorPointsConfig).
	Include(game.SetPokeDollarsConfig).
	Include(game.SwitchPartySlotsConfig).
	Include(game.SubmitChatMessageConfig).
	Include(game.SwitchChatChannelConfig).
	Include(game.SubmitChatCommandConfig).
	Include(game.SelectPlayerOptionConfig).
	Include(game.SetServerTimeConfig)

// messageCodec holds demarshallers and marshallers of messages.
var messageCodec = client.NewCodec().
	Join(loginCodec).
	Join(gameCodec)

// The main entry point to this game server application.
func main() {
	logger, err := createZapLogger()
	if err != nil {
		log.Fatal("Failed to create zap logger", err)
	}

	logger.Info("Starting PokeSync...")

	tcpHost := getTCPServerHostFromEnv()
	tcpPort := getTCPServerPortFromEnv()

	redisHost := getRedisHostFromEnv()
	redisPort := getRedisPortFromEnv()

	redisClient, err := connectToRedis(redisHost, redisPort)
	if err != nil {
		log.Fatal("Failed to connect to a Redis server instance", err)
	}

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

	routing := client.NewRouter(client.RouterConfig{
		PublicationTimeout: 1 * time.Second,
	})

	sessionConfig := game.SessionConfig{
		CommandLimit: 16,
		EventLimit:   256,
	}

	gameConfig := game.Config{
		IntervalRate:      50 * time.Millisecond,
		JobLimit:          game.Unbounded,
		WorkerCount:       runtime.NumCPU(),
		EntityLimit:       32768,
		ClockRate:         250 * time.Millisecond,
		ClockSynchronizer: game.NewGMT0Synchronizer(),
		Logger:            logger,
		SessionConfig:     sessionConfig,
	}

	statusConfig := status.Config{
		Logger:      logger,
		RefreshRate: 15 * time.Second,
	}

	discordConfig := discord.Config{}

	loginConfig := login.Config{
		JobLimit:          login.Unbounded,
		JobConsumeTimeout: time.Second,
		Logger:            logger,
		WorkerCount:       runtime.NumCPU(),
	}

	chatConfig := chat.Config{
		Logger: logger,
	}

	clientConfig := client.Config{
		Log:             logger,
		MessageCodec:    *messageCodec,
		ReadBufferSize:  512,
		WriteBufferSize: 2048,
		CommandLimit:    32,
	}

	serverConfig := server.Config{
		ClientConfig: clientConfig,
		Logger:       logger,
	}

	authenticator := login.NewAuthenticator(
		account.NewInMemoryRepository(),
		account.BasicPasswordMatcher(),
	)

	characters := character.NewInMemoryRepository()

	chatService := chat.NewService(chatConfig, routing)
	loginService := login.NewService(loginConfig, authenticator, routing)
	gameService := game.NewService(gameConfig, routing, characters, assetBundle)
	discordService := discord.NewService(discordConfig)
	statusService := status.NewService(statusConfig, status.NewRedisNotifier(redisClient), status.NewProvider(gameService))

	// should something go wrong and cause a panic, always safely
	// tear down these services
	defer func() {
		chatService.TearDown()
		loginService.TearDown()
		discordService.TearDown()
		gameService.TearDown()
		statusService.TearDown()
	}()

	logger.Info("Client build: ", ClientBuildNo)

	logger.Info("Game interval rate: ", gameConfig.IntervalRate)
	logger.Info("World entity limit: ", gameConfig.EntityLimit)
	logger.Info("Login worker count: ", loginConfig.WorkerCount)

	logger.Info("Item configs loaded: ", assetBundle.Items.Count())
	logger.Info("Npc configs loaded: ", assetBundle.Npcs.Count())
	logger.Info("Monster configs loaded: ", assetBundle.Monsters.Count())

	logger.Info("Upstream byte limit: ", clientConfig.ReadBufferSize)
	logger.Info("Downstream byte limit: ", clientConfig.WriteBufferSize)

	tcpListener := server.NewTcpListener(serverConfig, routing)
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
