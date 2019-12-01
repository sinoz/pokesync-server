package transport

import (
	"gitlab.com/pokesync/game-service/internal/game-service/client"
	"gitlab.com/pokesync/game-service/pkg/bytes"
)

var (
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
)

type MoveOrthographicCamera struct {
	X uint16
	Y uint16
}

type ResetOrthographicCamera struct{}

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
