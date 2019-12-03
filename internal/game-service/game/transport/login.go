package transport

import (
	"gitlab.com/pokesync/game-service/internal/game-service/client"
	"gitlab.com/pokesync/game-service/pkg/bytes"
)

var (
	UnableToFetchProfileConfig = client.MessageConfig{
		Kind:  client.UnableToFetchProfile,
		Topic: "unable_to_fetch_profile",
		New:   func() client.Message { return &UnableToFetchProfile{} },
	}

	RequestTimedOutConfig = client.MessageConfig{
		Kind:  client.LoginRequestTimedOut,
		Topic: "req_timeout",
		New:   func() client.Message { return &RequestTimedOut{} },
	}

	WorldFullConfig = client.MessageConfig{
		Kind:  client.WorldFull,
		Topic: "world_full",
		New:   func() client.Message { return &WorldFull{} },
	}

	LoginSuccessConfig = client.MessageConfig{
		Kind:  client.LoginSuccess,
		Topic: "login_success",
		New:   func() client.Message { return &LoginSuccess{} },
	}

	SelectCharacterConfig = client.MessageConfig{
		Kind:  client.SelectCharacter,
		Topic: "select_character",
		New:   func() client.Message { return &SelectCharacter{} },
	}

	CreateCharacterConfig = client.MessageConfig{
		Kind:  client.CreateCharacter,
		Topic: "create_character",
		New:   func() client.Message { return &CreateCharacter{} },
	}
)

type LoginSuccess struct {
	PID         uint16
	DisplayName string
	Gender      byte
	UserGroup   byte
	MapX        uint16
	MapZ        uint16
	LocalX      uint16
	LocalZ      uint16
}

type UnableToFetchProfile struct{}

type RequestTimedOut struct{}

type WorldFull struct{}

type CreateCharacter struct {
}

type SelectCharacter struct {
	Index byte
}

func (message *LoginSuccess) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.PID, _ = itr.ReadUInt16()
	message.DisplayName, _ = itr.ReadCString()
	message.Gender, _ = itr.ReadByte()
	message.UserGroup, _ = itr.ReadByte()

	message.MapX, _ = itr.ReadUInt16()
	message.MapZ, _ = itr.ReadUInt16()
	message.LocalX, _ = itr.ReadUInt16()
	message.LocalZ, _ = itr.ReadUInt16()
}

func (message *LoginSuccess) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.
		WriteInt16(int16(message.PID)).
		WriteCString(message.DisplayName).
		WriteByte(message.Gender).
		WriteByte(message.UserGroup).
		WriteInt16(int16(message.MapX)).
		WriteInt16(int16(message.MapZ)).
		WriteInt16(int16(message.LocalX)).
		WriteInt16(int16(message.LocalZ))

	return bldr.Build()
}

func (message *LoginSuccess) GetConfig() client.MessageConfig {
	return LoginSuccessConfig
}

func (r *UnableToFetchProfile) Demarshal(packet *client.Packet) {
}

func (r *UnableToFetchProfile) Marshal() *bytes.String {
	return bytes.EmptyString()
}

func (r *UnableToFetchProfile) GetConfig() client.MessageConfig {
	return UnableToFetchProfileConfig
}

func (r *RequestTimedOut) Demarshal(packet *client.Packet) {
}

func (r *RequestTimedOut) Marshal() *bytes.String {
	return bytes.EmptyString()
}

func (r *RequestTimedOut) GetConfig() client.MessageConfig {
	return RequestTimedOutConfig
}

func (r *WorldFull) Demarshal(packet *client.Packet) {
}

func (r *WorldFull) Marshal() *bytes.String {
	return bytes.EmptyString()
}

func (r *WorldFull) GetConfig() client.MessageConfig {
	return WorldFullConfig
}

func (message *CreateCharacter) Demarshal(packet *client.Packet) {
}

func (message *CreateCharacter) Marshal() *bytes.String {
	return bytes.EmptyString()
}

func (message *CreateCharacter) GetConfig() client.MessageConfig {
	return CreateCharacterConfig
}

func (message *SelectCharacter) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.Index, _ = itr.ReadByte()
}

func (message *SelectCharacter) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteByte(message.Index)

	return bldr.Build()
}

func (message *SelectCharacter) GetConfig() client.MessageConfig {
	return SelectCharacterConfig
}
