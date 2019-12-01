package transport

import (
	"gitlab.com/pokesync/game-service/internal/game-service/client"
	"gitlab.com/pokesync/game-service/pkg/bytes"
)

var (
	ChangeMovementTypeConfig = client.MessageConfig{
		Kind:  client.ChangeMoveType,
		Topic: "change_move_type",
		New:   func() client.Message { return &ChangeMovementType{} },
	}

	ClickTeleportConfig = client.MessageConfig{
		Kind:  client.ClickTeleport,
		Topic: "click_teleport",
		New:   func() client.Message { return &ClickTeleport{} },
	}

	AttachFollowerConfig = client.MessageConfig{
		Kind:  client.AttachFollower,
		Topic: "attach_follower",
		New:   func() client.Message { return &AttachFollower{} },
	}

	ClearFollowerConfig = client.MessageConfig{
		Kind:  client.ClearFollower,
		Topic: "clear_follower",
		New:   func() client.Message { return &ClearFollower{} },
	}

	EntityUpdateConfig = client.MessageConfig{
		Kind:  client.EntityUpdate,
		Topic: "entity_update",
		New:   func() client.Message { return &EntityUpdate{} },
	}

	FaceDirectionConfig = client.MessageConfig{
		Kind:  client.FaceDirection,
		Topic: "face_direction",
		New:   func() client.Message { return &FaceDirection{} },
	}

	InteractWithEntityConfig = client.MessageConfig{
		Kind:  client.InteractWithEntity,
		Topic: "interact_with_entity",
		New:   func() client.Message { return &InteractWithEntity{} },
	}

	MoveAvatarConfig = client.MessageConfig{
		Kind:  client.MoveAvatar,
		Topic: "move_avatar",
		New:   func() client.Message { return &MoveAvatar{} },
	}
)

type MoveAvatar struct {
	Direction byte
}

type ChangeMovementType struct {
	Type byte
}

type ClickTeleport struct {
	MapX   uint16
	MapZ   uint16
	LocalX uint16
	LocalZ uint16
}

type AttachFollower struct {
	PartySlot byte
}

type ClearFollower struct {
}

type FaceDirection struct {
	Direction byte
}

type InteractWithEntity struct {
	PID uint16
}

type EntityUpdate struct {
	// TODO
}

func (message *MoveAvatar) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.Direction, _ = itr.ReadByte()
}

func (message *MoveAvatar) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteByte(message.Direction)

	return bldr.Build()
}

func (message *MoveAvatar) GetConfig() client.MessageConfig {
	return MoveAvatarConfig
}

func (message *AttachFollower) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.PartySlot, _ = itr.ReadByte()
}

func (message *AttachFollower) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteByte(message.PartySlot)

	return bldr.Build()
}

func (message *AttachFollower) GetConfig() client.MessageConfig {
	return AttachFollowerConfig
}

func (message *ClearFollower) Demarshal(packet *client.Packet) {
}

func (message *ClearFollower) Marshal() *bytes.String {
	return bytes.EmptyString()
}

func (message *ClearFollower) GetConfig() client.MessageConfig {
	return ClearFollowerConfig
}

func (message *ChangeMovementType) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.Type, _ = itr.ReadByte()
}

func (message *ChangeMovementType) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteByte(message.Type)

	return bldr.Build()
}

func (message *ChangeMovementType) GetConfig() client.MessageConfig {
	return ChangeMovementTypeConfig
}

func (message *ClickTeleport) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.MapX, _ = itr.ReadUInt16()
	message.MapZ, _ = itr.ReadUInt16()
	message.LocalX, _ = itr.ReadUInt16()
	message.LocalZ, _ = itr.ReadUInt16()
}

func (message *ClickTeleport) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteInt16(int16(message.MapX))
	bldr.WriteInt16(int16(message.MapZ))
	bldr.WriteInt16(int16(message.LocalX))
	bldr.WriteInt16(int16(message.LocalZ))

	return bldr.Build()
}

func (message *ClickTeleport) GetConfig() client.MessageConfig {
	return ClickTeleportConfig
}

func (message *FaceDirection) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.Direction, _ = itr.ReadByte()
}

func (message *FaceDirection) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteByte(message.Direction)

	return bldr.Build()
}

func (message *FaceDirection) GetConfig() client.MessageConfig {
	return FaceDirectionConfig
}

func (message *EntityUpdate) Demarshal(packet *client.Packet) {
	// TODO
}

func (message *EntityUpdate) Marshal() *bytes.String {
	return bytes.EmptyString() // TODO
}

func (message *EntityUpdate) GetConfig() client.MessageConfig {
	return EntityUpdateConfig
}

func (message *InteractWithEntity) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.PID, _ = itr.ReadUInt16()
}

func (message *InteractWithEntity) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteInt16(int16(message.PID))

	return bldr.Build()
}

func (message *InteractWithEntity) GetConfig() client.MessageConfig {
	return InteractWithEntityConfig
}
