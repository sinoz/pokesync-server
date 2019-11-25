package chat

import (
	"gitlab.com/pokesync/game-service/internal/pokesync/character"
)

// ChannelId is an alias of an int that represents the channel id.
type ChannelId int

// MessageText is an alias of the text contents of a chat message.
type MessageText string

// Message is a chat message broadcasted by a public chat participant.
type Message struct {
	sender    character.DisplayName
	channelId ChannelId
	text      MessageText
}
