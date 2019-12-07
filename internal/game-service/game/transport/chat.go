package transport

import (
	"gitlab.com/pokesync/game-service/internal/game-service/client"
	"gitlab.com/pokesync/game-service/pkg/bytes"
)

var (
	SubmitChatCommandConfig = client.MessageConfig{
		Kind:  client.SubmitChatCmd,
		Topic: "submit_chat_cmd",
		New:   func() client.Message { return &SubmitChatCommand{} },
	}
)

type SubmitChatCommand struct {
	Trigger   string
	Arguments []string
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
