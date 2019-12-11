package transport

import (
	"gitlab.com/pokesync/game-service/internal/game-service/client"
	"gitlab.com/pokesync/game-service/pkg/bytes"
)

var (
	SwitchPartySlotsConfig = client.MessageConfig{
		Kind:  client.SwitchPartySlots,
		Topic: "switch_party_slots",
		New:   func() client.Message { return &SetPartySlot{} },
	}

	SetPartySlotConfig = client.MessageConfig{
		Kind:  client.SetPartySlot,
		Topic: "set_party_slot",
		New:   func() client.Message { return &SetPartySlot{} },
	}
)

type SwitchPartySlots struct {
	SlotFrom byte
	SlotTo   byte
}

type SetPartySlot struct {
	Slot            byte
	MonsterID       uint16
	Gender          byte
	Coloration      byte
	StatusCondition byte
}

func (message *SwitchPartySlots) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.SlotFrom, _ = itr.ReadByte()
	message.SlotTo, _ = itr.ReadByte()
}

func (message *SwitchPartySlots) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteByte(message.SlotFrom)
	bldr.WriteByte(message.SlotTo)

	return bldr.Build()
}

func (message *SwitchPartySlots) GetConfig() client.MessageConfig {
	return SwitchPartySlotsConfig
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
