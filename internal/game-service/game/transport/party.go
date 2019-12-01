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

	SwitchPartySlotsConfig = client.MessageConfig{
		Kind:  client.SwitchPartySlots,
		Topic: "switch_party_slots",
		New:   func() client.Message { return &SwitchPartySlots{} },
	}
)

type SetPartySlot struct {
	Slot byte
}

type SwitchPartySlots struct {
	From byte
	To   byte
}

func (message *SetPartySlot) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.Slot, _ = itr.ReadByte()
}

func (message *SetPartySlot) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteByte(message.Slot)

	return bldr.Build()
}

func (message *SetPartySlot) GetConfig() client.MessageConfig {
	return SetPartySlotConfig
}

func (message *SwitchPartySlots) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.From, _ = itr.ReadByte()
	message.To, _ = itr.ReadByte()
}

func (message *SwitchPartySlots) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteByte(message.From)
	bldr.WriteByte(message.To)

	return bldr.Build()
}

func (message *SwitchPartySlots) GetConfig() client.MessageConfig {
	return SwitchPartySlotsConfig
}
