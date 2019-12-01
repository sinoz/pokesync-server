package transport

import (
	"gitlab.com/pokesync/game-service/internal/game-service/client"
	"gitlab.com/pokesync/game-service/pkg/bytes"
)

var (
	CloseDialogueConfig = client.MessageConfig{
		Kind:  client.CloseDialogue,
		Topic: "close_dialogue",
		New:   func() client.Message { return &CloseDialogue{} },
	}

	ContinueDialogueConfig = client.MessageConfig{
		Kind:  client.ContinueDialogue,
		Topic: "continue_dialog",
		New:   func() client.Message { return &ContinueDialogue{} },
	}

	SetDonatorPointsConfig = client.MessageConfig{
		Kind:  client.SetDonatorPoints,
		Topic: "set_donator_pts",
		New:   func() client.Message { return &SetDonatorPoints{} },
	}

	SetPokeDollarsConfig = client.MessageConfig{
		Kind:  client.SetPokeDollar,
		Topic: "set_pokedollar",
		New:   func() client.Message { return &SetPokeDollar{} },
	}

	SetServerTimeConfig = client.MessageConfig{
		Kind:  client.SetServerTime,
		Topic: "set_server_time",
		New:   func() client.Message { return &SetServerTime{} },
	}

	SelectPlayerOptionConfig = client.MessageConfig{
		Kind:  client.SelectPlayerOpt,
		Topic: "select_plr_opt",
		New:   func() client.Message { return &SelectPlayerOption{} },
	}
)

type SetPokeDollar struct {
	Amount uint32
}

type SetDonatorPoints struct {
	Amount uint32
}

type SetServerTime struct {
	Hour   byte
	Minute byte
}

type SelectPlayerOption struct {
	PID    uint16
	Option byte
}

type ContinueDialogue struct {
}

type CloseDialogue struct {
}

func (message *ContinueDialogue) Demarshal(packet *client.Packet) {
}

func (message *ContinueDialogue) Marshal() *bytes.String {
	return bytes.EmptyString()
}

func (message *ContinueDialogue) GetConfig() client.MessageConfig {
	return ContinueDialogueConfig
}

func (message *CloseDialogue) Demarshal(packet *client.Packet) {
}

func (message *CloseDialogue) Marshal() *bytes.String {
	return bytes.EmptyString()
}

func (message *CloseDialogue) GetConfig() client.MessageConfig {
	return CloseDialogueConfig
}

func (message *SetPokeDollar) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.Amount, _ = itr.ReadUInt32()
}

func (message *SetPokeDollar) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteInt32(int32(message.Amount))

	return bldr.Build()
}

func (message *SetPokeDollar) GetConfig() client.MessageConfig {
	return SetPokeDollarsConfig
}

func (message *SetDonatorPoints) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.Amount, _ = itr.ReadUInt32()
}

func (message *SetDonatorPoints) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteInt32(int32(message.Amount))

	return bldr.Build()
}

func (message *SetDonatorPoints) GetConfig() client.MessageConfig {
	return SetDonatorPointsConfig
}

func (message *SelectPlayerOption) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.PID, _ = itr.ReadUInt16()
	message.Option, _ = itr.ReadByte()
}

func (message *SelectPlayerOption) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteInt16(int16(message.PID))
	bldr.WriteByte(message.Option)

	return bldr.Build()
}

func (message *SelectPlayerOption) GetConfig() client.MessageConfig {
	return SelectPlayerOptionConfig
}

func (message *SetServerTime) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.Hour, _ = itr.ReadByte()
	message.Minute, _ = itr.ReadByte()
}

func (message *SetServerTime) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteByte(message.Hour)
	bldr.WriteByte(message.Minute)

	return bldr.Build()
}

func (message *SetServerTime) GetConfig() client.MessageConfig {
	return SetServerTimeConfig
}
