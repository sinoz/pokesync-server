package chat

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
