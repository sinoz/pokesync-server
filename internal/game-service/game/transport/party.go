package transport

import (
	"gitlab.com/pokesync/game-service/internal/game-service/client"
	"gitlab.com/pokesync/game-service/pkg/bytes"
)

var (
	SetPartySlotConfig = client.MessageConfig{
		Kind:  client.SetPartySlot,
		Topic: "set_party_slot",
		New:   func() client.Message { return &SetPartySlot{} },
	}
)

type SetPartySlot struct {
	Slot            byte
	MonsterID       uint16
	Gender          byte
	Coloration      byte
	StatusCondition byte
}

func (message *SetPartySlot) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.Slot, _ = itr.ReadByte()
	message.MonsterID, _ = itr.ReadUInt16()
	message.Gender, _ = itr.ReadByte()
	message.Coloration, _ = itr.ReadByte()
	message.StatusCondition, _ = itr.ReadByte()
}

func (message *SetPartySlot) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteByte(message.Slot)
	bldr.WriteInt16(int16(message.MonsterID))
	bldr.WriteByte(message.Gender)
	bldr.WriteByte(message.Coloration)
	bldr.WriteByte(message.StatusCondition)

	return bldr.Build()
}

func (message *SetPartySlot) GetConfig() client.MessageConfig {
	return SetPartySlotConfig
}
