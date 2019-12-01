package transport

import (
	"gitlab.com/pokesync/game-service/internal/game-service/client"
	"gitlab.com/pokesync/game-service/pkg/bytes"
)

var (
	DisplayChatMessageConfig = client.MessageConfig{
		Kind:  client.DisplayChatMsg,
		Topic: "display_chat_msg",
		New:   func() client.Message { return &DisplayChatMessage{} },
	}

	SubmitChatCommandConfig = client.MessageConfig{
		Kind:  client.SubmitChatCmd,
		Topic: "submit_chat_cmd",
		New:   func() client.Message { return &SubmitChatCommand{} },
	}

	SwitchChatChannelConfig = client.MessageConfig{
		Kind:  client.SwitchChatChannel,
		Topic: "switch_chat_channel",
		New:   func() client.Message { return &SwitchChatChannel{} },
	}

	SubmitChatMessageConfig = client.MessageConfig{
		Kind:  client.SubmitChatMsg,
		Topic: "submit_chat_msg",
		New:   func() client.Message { return &SubmitChatMessage{} },
	}

	SelectChatChannelConfig = client.MessageConfig{
		Kind:  client.SelectChatChannel,
		Topic: "select_chat_channel",
		New:   func() client.Message { return &SelectChatChannel{} },
	}
)

type SubmitChatCommand struct {
	Trigger   string
	Arguments []string
}

type SwitchChatChannel struct {
	ChannelId byte
}

type SubmitChatMessage struct {
	Text string
}

type DisplayChatMessage struct {
	ChannelId   byte
	DisplayName string
	Text        string
}

type SelectChatChannel struct {
	ChannelId byte
}

func (message *DisplayChatMessage) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.ChannelId, _ = itr.ReadByte()
	message.DisplayName, _ = itr.ReadCString()
	message.Text, _ = itr.ReadCString()
}

func (message *DisplayChatMessage) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteByte(message.ChannelId)
	bldr.WriteCString(message.DisplayName)
	bldr.WriteCString(message.Text)

	return bldr.Build()
}

func (message *DisplayChatMessage) GetConfig() client.MessageConfig {
	return DisplayChatMessageConfig
}

func (message *SwitchChatChannel) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.ChannelId, _ = itr.ReadByte()
}

func (message *SwitchChatChannel) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteByte(message.ChannelId)

	return bldr.Build()
}

func (message *SwitchChatChannel) GetConfig() client.MessageConfig {
	return SwitchChatChannelConfig
}

func (message *SubmitChatMessage) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.Text, _ = itr.ReadCString()
}

func (message *SubmitChatMessage) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteCString(message.Text)

	return bldr.Build()
}

func (message *SubmitChatMessage) GetConfig() client.MessageConfig {
	return SubmitChatMessageConfig
}

func (message *SubmitChatCommand) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.Trigger, _ = itr.ReadCString()
	for i := 0; itr.IsReadable(); i++ {
		message.Arguments[i], _ = itr.ReadCString()
	}
}

func (message *SubmitChatCommand) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteCString(message.Trigger)
	for _, argument := range message.Arguments {
		bldr.WriteCString(argument)
	}

	return bldr.Build()
}

func (message *SubmitChatCommand) GetConfig() client.MessageConfig {
	return SubmitChatCommandConfig
}

func (message *SelectChatChannel) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.ChannelId, _ = itr.ReadByte()
}

func (message *SelectChatChannel) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteByte(message.ChannelId)

	return bldr.Build()
}

func (message *SelectChatChannel) GetConfig() client.MessageConfig {
	return SelectChatChannelConfig
}
