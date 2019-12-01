package transport

import (
	"gitlab.com/pokesync/game-service/internal/game-service/client"
	"gitlab.com/pokesync/game-service/pkg/bytes"
)

var (
	RefreshMapConfig = client.MessageConfig{
		Kind:  client.MapRefresh,
		Topic: "move_avatar",
		New:   func() client.Message { return &RefreshMap{} },
	}
)

type RefreshMap struct {
	MapX uint16
	MapZ uint16
}

func (message *RefreshMap) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.MapX, _ = itr.ReadUInt16()
	message.MapZ, _ = itr.ReadUInt16()
}

func (message *RefreshMap) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteInt16(int16(message.MapX))
	bldr.WriteInt16(int16(message.MapZ))

	return bldr.Build()
}

func (message *RefreshMap) GetConfig() client.MessageConfig {
	return RefreshMapConfig
}
