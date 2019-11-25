package game

import (
	"gitlab.com/pokesync/game-service/internal/game-service/client"
	"gitlab.com/pokesync/game-service/pkg/bytes"
)

var (
	UnableToFetchProfileConfig = client.MessageConfig{
		Kind: client.UnableToFetchProfile,
		New:  func() client.Message { return &UnableToFetchProfile{} },
	}

	LoginSuccessConfig = client.MessageConfig{
		Kind:  client.LoginSuccess,
		Topic: "login_success",
		New:   func() client.Message { return &LoginSuccess{} },
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

	DisplayChatMessageConfig = client.MessageConfig{
		Kind:  client.DisplayChatMsg,
		Topic: "display_chat_msg",
		New:   func() client.Message { return &DisplayChatMessage{} },
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

	RefreshMapConfig = client.MessageConfig{
		Kind:  client.MapRefresh,
		Topic: "move_avatar",
		New:   func() client.Message { return &RefreshMap{} },
	}

	MoveCameraConfig = client.MessageConfig{
		Kind:  client.MoveOrthoCamera,
		Topic: "move_camera",
		New:   func() client.Message { return &MoveOrthographicCamera{} },
	}

	ResetCameraConfig = client.MessageConfig{
		Kind:  client.ResetOrthoCamera,
		Topic: "reset_camera",
		New:   func() client.Message { return &ResetOrthographicCamera{} },
	}

	MoveAvatarConfig = client.MessageConfig{
		Kind:  client.MoveAvatar,
		Topic: "move_avatar",
		New:   func() client.Message { return &MoveAvatar{} },
	}

	SelectChatChannelConfig = client.MessageConfig{
		Kind:  client.SelectChatChannel,
		Topic: "select_chat_channel",
		New:   func() client.Message { return &SelectChatChannel{} },
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

type MoveAvatar struct {
	Direction byte
}

type RefreshMap struct {
	MapX uint16
	MapZ uint16
}

type MoveOrthographicCamera struct {
	X uint16
	Y uint16
}

type ResetOrthographicCamera struct{}

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

type CreateCharacter struct {
}

type SelectCharacter struct {
	Index byte
}

type EntityUpdate struct {
	// TODO
}

type SelectChatChannel struct {
	ChannelId byte
}

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

type SetPartySlot struct {
	Slot byte
}

type SwitchPartySlots struct {
	From byte
	To   byte
}

type ContinueDialogue struct {
}

type CloseDialogue struct {
}

type DisplayChatMessage struct {
	ChannelId   byte
	DisplayName string
	Text        string
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

func (r *UnableToFetchProfile) Demarshal(packet *client.Packet) {
}

func (r *UnableToFetchProfile) Marshal() *bytes.String {
	return bytes.EmptyString()
}

func (r *UnableToFetchProfile) GetConfig() client.MessageConfig {
	return UnableToFetchProfileConfig
}

func (message *LoginSuccess) GetConfig() client.MessageConfig {
	return LoginSuccessConfig
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

func (message *ResetOrthographicCamera) Demarshal(packet *client.Packet) {
}

func (message *ResetOrthographicCamera) Marshal() *bytes.String {
	return bytes.EmptyString()
}

func (message *ResetOrthographicCamera) GetConfig() client.MessageConfig {
	return ResetCameraConfig
}

func (message *MoveOrthographicCamera) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	message.X, _ = itr.ReadUInt16()
	message.Y, _ = itr.ReadUInt16()
}

func (message *MoveOrthographicCamera) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.WriteInt16(int16(message.X))
	bldr.WriteInt16(int16(message.Y))

	return bldr.Build()
}

func (message *MoveOrthographicCamera) GetConfig() client.MessageConfig {
	return MoveCameraConfig
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
