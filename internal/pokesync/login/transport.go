package login

import (
	"gitlab.com/pokesync/game-service/internal/pokesync/client"
	"gitlab.com/pokesync/game-service/pkg/bytes"
)

var (
	RequestConfig = client.MessageConfig{
		Kind:  client.RequestLogin,
		Topic: "login_request",
		New:   func() client.Message { return &Request{} },
	}

	InvalidCredentialsConfig = client.MessageConfig{
		Kind: client.InvalidCredentials,
		New:  func() client.Message { return &InvalidCredentials{} },
	}

	AlreadyLoggedInConfig = client.MessageConfig{
		Kind: client.AlreadyLoggedIn,
		New:  func() client.Message { return &AlreadyLoggedIn{} },
	}

	WorldFullConfig = client.MessageConfig{
		Kind: client.WorldFull,
		New:  func() client.Message { return &WorldFull{} },
	}

	AccountDisabledConfig = client.MessageConfig{
		Kind: client.AccountDisabled,
		New:  func() client.Message { return &AccountDisabled{} },
	}

	ErrorDuringAccountFetchConfig = client.MessageConfig{
		Kind: client.UnableToFetchProfile,
		New:  func() client.Message { return &ErrorDuringAccountFetch{} },
	}

	RequestTimedOutConfig = client.MessageConfig{
		Kind: client.LoginRequestTimedOut,
		New:  func() client.Message { return &RequestTimedOut{} },
	}
)

type Request struct {
	MajorVersion uint8
	MinorVersion uint8
	PatchVersion uint8

	Email    string
	Password string
}

type InvalidCredentials struct{}

type WorldFull struct{}

type AlreadyLoggedIn struct{}

type AccountDisabled struct{}

type ErrorDuringAccountFetch struct{}

type RequestTimedOut struct{}

func (r *Request) Demarshal(packet *client.Packet) {
	itr := packet.Bytes.Iterator()

	r.MajorVersion, _ = itr.ReadByte()
	r.MinorVersion, _ = itr.ReadByte()
	r.PatchVersion, _ = itr.ReadByte()

	r.Email, _ = itr.ReadCString()
	r.Password, _ = itr.ReadCString()
}

func (r *Request) Marshal() *bytes.String {
	bldr := bytes.NewDefaultBuilder()

	bldr.
		WriteByte(r.MajorVersion).
		WriteByte(r.MinorVersion).
		WriteByte(r.PatchVersion).
		WriteCString(r.Email).
		WriteCString(r.Password)

	return bldr.Build()
}

func (r *Request) GetConfig() client.MessageConfig {
	return RequestConfig
}

func (r *AccountDisabled) Demarshal(packet *client.Packet) {
}

func (r *AccountDisabled) Marshal() *bytes.String {
	return bytes.EmptyString()
}

func (r *AccountDisabled) GetConfig() client.MessageConfig {
	return AccountDisabledConfig
}

func (r *WorldFull) Demarshal(packet *client.Packet) {
}

func (r *WorldFull) Marshal() *bytes.String {
	return bytes.EmptyString()
}

func (r *WorldFull) GetConfig() client.MessageConfig {
	return WorldFullConfig
}

func (r *AlreadyLoggedIn) Demarshal(packet *client.Packet) {
}

func (r *AlreadyLoggedIn) Marshal() *bytes.String {
	return bytes.EmptyString()
}

func (r *AlreadyLoggedIn) GetConfig() client.MessageConfig {
	return AlreadyLoggedInConfig
}

func (r *InvalidCredentials) Demarshal(packet *client.Packet) {
}

func (r *InvalidCredentials) Marshal() *bytes.String {
	return bytes.EmptyString()
}

func (r *InvalidCredentials) GetConfig() client.MessageConfig {
	return InvalidCredentialsConfig
}

func (r *RequestTimedOut) Demarshal(packet *client.Packet) {
}

func (r *RequestTimedOut) Marshal() *bytes.String {
	return bytes.EmptyString()
}

func (r *RequestTimedOut) GetConfig() client.MessageConfig {
	return RequestTimedOutConfig
}

func (r *ErrorDuringAccountFetch) Demarshal(packet *client.Packet) {
}

func (r *ErrorDuringAccountFetch) Marshal() *bytes.String {
	return bytes.EmptyString()
}

func (r *ErrorDuringAccountFetch) GetConfig() client.MessageConfig {
	return ErrorDuringAccountFetchConfig
}
