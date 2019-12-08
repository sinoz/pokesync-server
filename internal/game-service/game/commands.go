package game

import "gitlab.com/pokesync/game-service/internal/game-service/game/entity"

// ChatCommandHandler handles a chat command with the given arguments.
// May return an error which is to flow upwards through the call chain.
type ChatCommandHandler func(entity *entity.Entity, arguments []string) error

// ChatCommandRegistry is a registry of ChatCommandHandler's.
type ChatCommandRegistry struct {
	chatCommands map[string]ChatCommandHandler
}

// NewChatCommandRegistry constructs a new instance of a ChatCommandRegistry.
func NewChatCommandRegistry() *ChatCommandRegistry {
	return &ChatCommandRegistry{chatCommands: make(map[string]ChatCommandHandler)}
}

// Put inserts the given ChatCommandHandler into the registry under the
// specified trigger.
func (registry *ChatCommandRegistry) Put(trigger string, handler ChatCommandHandler) {
	registry.chatCommands[trigger] = handler
}

// Remove removes any ChatCommandHandler that is associated with the specified
// trigger.
func (registry *ChatCommandRegistry) Remove(trigger string) {
	delete(registry.chatCommands, trigger)
}

// Get looks up a ChatCommandHandler by its specified trigger.
func (registry *ChatCommandRegistry) Get(trigger string) (ChatCommandHandler, bool) {
	handler, exists := registry.chatCommands[trigger]
	return handler, exists
}

// submitChatCommand is a submitChatCommandHandler that looks-up and calls
// the ChatCommandHandler that is associated with a given trigger.
func submitChatCommand(registry *ChatCommandRegistry) submitChatCommandHandler {
	return func(entity *entity.Entity, trigger string, arguments []string) error {
		handle, exists := registry.Get(trigger)
		if !exists {
			return nil
		}

		return handle(entity, arguments)
	}
}
