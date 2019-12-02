package client

import (
	"bufio"
	"errors"

	"gitlab.com/pokesync/game-service/pkg/bytes"
)

const (
	// Client -> Server
	RequestLogin       PacketKind = 0
	CreateCharacter    PacketKind = 1
	SelectCharacter    PacketKind = 2
	AttachFollower     PacketKind = 3
	ClearFollower      PacketKind = 4
	ChangeMoveType     PacketKind = 5
	ContinueDialogue   PacketKind = 6
	FaceDirection      PacketKind = 7
	InteractWithEntity PacketKind = 8
	MoveAvatar         PacketKind = 9
	SelectPlayerOpt    PacketKind = 10
	SwitchPartySlots   PacketKind = 11
	SubmitChatMsg      PacketKind = 12
	SubmitChatCmd      PacketKind = 13
	ClickTeleport      PacketKind = 14
	SelectChatChannel  PacketKind = 15

	// Server -> Client
	MapRefresh           PacketKind = 238
	SwitchChatChannel    PacketKind = 239
	SetServerTime        PacketKind = 240
	SetPokeDollar        PacketKind = 241
	SetPartySlot         PacketKind = 242
	SetDonatorPoints     PacketKind = 243
	ResetOrthoCamera     PacketKind = 244
	MoveOrthoCamera      PacketKind = 245
	EntityUpdate         PacketKind = 246
	CloseDialogue        PacketKind = 247
	DisplayChatMsg       PacketKind = 248
	WorldFull            PacketKind = 249
	LoginRequestTimedOut PacketKind = 250
	UnableToFetchProfile PacketKind = 251
	InvalidCredentials   PacketKind = 252
	AlreadyLoggedIn      PacketKind = 253
	AccountDisabled      PacketKind = 254
	LoginSuccess         PacketKind = 255
)

// PacketKind is a kind of packet.
type PacketKind uint8

// Packet is a structured unit of data that can be transferred
// across the wire.
type Packet struct {
	Kind  PacketKind
	Bytes *bytes.String
}

// MessageConfig holds configurations specifically for Message's.
type MessageConfig struct {
	Kind  PacketKind
	Topic Topic
	New   func() Message
}

// Demarshaller demarshals bytes of a Packet into a message.
type Demarshaller interface {
	Demarshal(packet *Packet)
}

// Marshaller marshals a message into bytes.
type Marshaller interface {
	Marshal() *bytes.String
}

// Message represents the abstraction for a structure of data
// sent across the wire.
type Message interface {
	Demarshaller
	Marshaller

	GetConfig() MessageConfig
}

// Codec holds message configs for all messages.
type Codec struct {
	configs map[PacketKind]MessageConfig
}

// NewCodec constructs a new message Codec.
func NewCodec() *Codec {
	return &Codec{configs: make(map[PacketKind]MessageConfig)}
}

// HeaderLength returns the length of the Packet's header, in bytes.
func (p *Packet) HeaderLength() int {
	return 1 + ComputeRawVarInt32Size(p.PayloadLength())
}

// PayloadLength returns the length of the Packet's payload, in bytes.
func (p *Packet) PayloadLength() int {
	return p.Bytes.Length()
}

// TotalLength returns the complete length of the Packet, in bytes.
func (p *Packet) TotalLength() int {
	return p.HeaderLength() + p.PayloadLength()
}

// ForkPacket forks a Packet from the given bufio.Reader. May return an
// error if something went wrong whilst reading (like when the underlying
// connection has been closed).
func ForkPacket(reader *bufio.Reader) (*Packet, error) {
	kind, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	length, err := readRawVarInt32(reader)
	if err != nil {
		return nil, err
	}

	payload := make([]byte, length)
	if _, err := reader.Read(payload); err != nil {
		return nil, err
	}

	return &Packet{Kind: PacketKind(kind), Bytes: bytes.StringWrap(payload)}, nil
}

// JoinPacket joins the given Packet into a byte stream of the given
// bufio.Writer. May return an error if something went wrong whilst writing.
func JoinPacket(writer *bufio.Writer, packet *Packet) error {
	var err error

	if err = writer.WriteByte(byte(packet.Kind)); err != nil {
		return err
	}

	if err := writeRawVarInt32(writer, packet.Bytes.Length()); err != nil {
		return err
	}

	if _, err := writer.Write(packet.Bytes.ToByteArray()); err != nil {
		return err
	}

	return nil
}

// readRawVarInt32 reads a value as a 32-bit variable integer.
func readRawVarInt32(reader *bufio.Reader) (int, error) {
	v, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	tmp := int8(v)
	if tmp >= 0 {
		return int(tmp), nil
	}

	result := int(tmp & 127)

	v, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	tmp = int8(v)
	if tmp >= 0 {
		result |= int(tmp) << 7
	} else {
		result |= (int(tmp) & 127) << 7

		v, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}

		tmp = int8(v)
		if tmp >= 0 {
			result |= int(tmp) << 14
		} else {
			result |= (int(tmp) & 127) << 14

			v, err = reader.ReadByte()
			if err != nil {
				return 0, err
			}

			tmp = int8(v)
			if tmp >= 0 {
				result |= int(tmp) << 21
			} else {
				result |= (int(tmp) & 127) << 21

				v, err = reader.ReadByte()
				if err != nil {
					return 0, err
				}

				tmp = int8(v)
				result |= int(tmp) << 28
				if tmp < 0 {
					return 0, errors.New("malformed varint")
				}
			}
		}
	}

	return result, nil
}

// writeRawVarInt32 writes the given value as a 32-bit variable integer,
// which actual size may vary depending on the amount of bytes the value
// is to occupy in the byte stream.
func writeRawVarInt32(writer *bufio.Writer, value int) error {
	for {
		if (value & ^0x7) == 0 {
			return writer.WriteByte(byte(value))
		}

		if err := writer.WriteByte(byte((value & 0x7F) | 0x80)); err != nil {
			return err
		}

		value = int(uint(value) >> 7)
	}
}

// ComputeRawVarInt32Size computes the size of a length field in bytes,
// depending on the given value of the length field.
func ComputeRawVarInt32Size(value int) int {
	if (value & (0xffffffff << 7)) == 0 {
		return 1
	}

	if (value & (0xffffffff << 14)) == 0 {
		return 2
	}

	if (value & (0xffffffff << 21)) == 0 {
		return 3
	}

	if (value & (0xffffffff << 28)) == 0 {
		return 4
	}

	return 5
}

// Join joins the given Codec with this Codec.
func (codec *Codec) Join(other *Codec) *Codec {
	for key, value := range other.configs {
		codec.configs[key] = value
	}

	return codec
}

// Include includes the given MessageConfig to deal with a specific kind
// of message.
func (codec *Codec) Include(config MessageConfig) *Codec {
	codec.configs[config.Kind] = config
	return codec
}

// GetConfig looks up a MessageConfig by the given PacketKind.
func (codec *Codec) GetConfig(kind PacketKind) (MessageConfig, bool) {
	config, exists := codec.configs[kind]
	return config, exists
}
