package main

import (
	"gitlab.com/pokesync/game-service/internal/pokesync/account"
	"gitlab.com/pokesync/game-service/internal/pokesync/character"
	"gitlab.com/pokesync/game-service/internal/pokesync/chat"
	"gitlab.com/pokesync/game-service/internal/pokesync/client"
	"gitlab.com/pokesync/game-service/internal/pokesync/game"
	"gitlab.com/pokesync/game-service/internal/pokesync/login"
	"gitlab.com/pokesync/game-service/internal/pokesync/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

// ClientBuildNo is the build number of the client this game
// server is supporting.
const ClientBuildNo = client.BuildNumber(1)

// ServerHostEnv is the name of the environment variable of
// the server host.
const ServerHostEnv = "POKESYNC_HOST"

// ServerPortEnv is the name of the environment variable of
// the server port.
const ServerPortEnv = "POKESYNC_PORT"

// DefaultServerHost is the port to fallback to if no environment
// variable is set.
const DefaultServerHost = "localhost"

// DefaultServerPort is the port to fallback to if no environment
// variable is set.
const DefaultServerPort = 23192

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

	gameConfig := game.Config{
		IntervalRate: 50 * time.Millisecond,
		JobLimit:     game.Unbounded,
		EntityLimit:  32768,
		WorkerCount:  runtime.NumCPU(),
	}

	loginConfig := login.Config{
		JobLimit:          login.Unbounded,
		JobConsumeTimeout: time.Second,
		WorkerCount:       runtime.NumCPU(),
	}

	chatConfig := chat.Config{
		Logger: logger,
	}

	clientConfig := client.Config{
		Log:          logger,
		MessageCodec: *messageCodec,

		ReadBufferSize:  512,
		WriteBufferSize: 2048,

		CommandLimit: 32,
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

	chat.NewService(chatConfig, routing)
	login.NewService(loginConfig, authenticator, routing)
	game.NewService(gameConfig, routing, characters, assetBundle)

	host := getServerHostFromEnv()
	port := getServerPortFromEnv()

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
	if err := tcpListener.Bind(host, port); err != nil {
		logger.Fatal(err)
	}
}

// createZapLogger constructs a sugared zap logger.
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

// getServerHostFromEnv extracts a host from the user's environment
// variables. If no environment variable is set, a fallback value is returned.
func getServerHostFromEnv() string {
	host := os.Getenv(ServerHostEnv)
	if len(host) == 0 {
		return DefaultServerHost
	}

	return host
}

// getServerPortFromEnv extracts a port number from the user's environment
// variables. If no environment variable is set, a fallback value is returned.
func getServerPortFromEnv() int {
	port, err := strconv.Atoi(os.Getenv(ServerPortEnv))
	if err != nil {
		return DefaultServerPort
	}

	return port
}
