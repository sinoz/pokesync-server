package game

import (
	"errors"

	"gitlab.com/pokesync/game-service/internal/game-service/game/session"
	"gitlab.com/pokesync/game-service/internal/game-service/game/transport"
)

const (
	// PartyBeltCapacity is the absolute capacity of a player's
	// monster party belt.
	PartyBeltCapacity = 6
)

// PartyBeltUpdateListener listens for changes made to the PartyBelt.
type PartyBeltUpdateListener interface {
	Updated(slot int, monster Monster)
}

// PartyBeltSessionListener is a PartyBeltUpdateListener that listens
// for changes made to the PartyBelt to visually apply these
// changes to the Client as well.
type PartyBeltSessionListener struct {
	session *session.Session
}

// PartyBelt is the belt of monsters a player carries with him/her.
type PartyBelt struct {
	monsters []*Monster
	Size     int

	listeners []PartyBeltUpdateListener
}

// NewPartyBelt constructs a new instance of a monster PartyBelt.
func NewPartyBelt() *PartyBelt {
	return &PartyBelt{}
}

func (belt *PartyBelt) Add(monster *Monster) bool {
	if belt.IsFull() {
		return false
	}

	belt.monsters = append(belt.monsters, monster)
	belt.Size++

	belt.notifySlotUpdated(belt.Size-1, *monster)

	return true
}

func (belt *PartyBelt) Swap(slotFrom, slotTo int) error {
	if err := belt.checkBoundaries(slotFrom); err != nil {
		return err
	}

	if err := belt.checkBoundaries(slotTo); err != nil {
		return err
	}

	monsterA := belt.monsters[slotFrom]
	monsterB := belt.monsters[slotTo]

	belt.monsters[slotFrom] = monsterB
	belt.monsters[slotTo] = monsterA

	belt.notifySlotUpdated(slotFrom, *monsterB)
	belt.notifySlotUpdated(slotTo, *monsterA)

	return nil
}

func (belt *PartyBelt) Clear(slot int) (*Monster, error) {
	if err := belt.checkBoundaries(slot); err != nil {
		return nil, err
	}

	before := belt.monsters[slot]
	belt.Size--

	belt.notifySlotUpdated(slot, *before)
	return before, nil
}

func (belt *PartyBelt) Set(slot int, monster *Monster) (*Monster, error) {
	if err := belt.checkBoundaries(slot); err != nil {
		return nil, err
	}

	before := belt.monsters[slot]
	belt.monsters[slot] = monster
	belt.Size++

	belt.notifySlotUpdated(slot, *monster)
	return before, nil
}

func (belt *PartyBelt) Get(slot int) (*Monster, error) {
	if err := belt.checkBoundaries(slot); err != nil {
		return nil, err
	}

	return belt.monsters[slot], nil
}

// IsEmpty returns whether there are any monsters in this belt.
func (belt *PartyBelt) IsEmpty() bool {
	return belt.Size == 0
}

// IsFull returns whether this belt is full or not.
func (belt *PartyBelt) IsFull() bool {
	return belt.Size == PartyBeltCapacity
}

// checkBoundaries checks if the given slot is valid. If not,
// it returns an error.
func (belt *PartyBelt) checkBoundaries(slot int) error {
	if slot < 0 {
		return errors.New("given slot is negative")
	}

	if slot >= belt.Size {
		return errors.New("given slot exceeds current belt size")
	}

	return nil
}

// notifySlotUpdated notifies every subscribed listener
// of the PartyBelt's updated slot.
func (belt *PartyBelt) notifySlotUpdated(slot int, monster Monster) {
	for _, listener := range belt.listeners {
		listener.Updated(slot, monster)
	}
}

// AddListener subscribes the given listener to receive notifications.
func (belt *PartyBelt) AddListener(listener PartyBeltUpdateListener) {
	belt.listeners = append(belt.listeners, listener)
}

// RemoveListener unsubscribes the given listener from receiving notifications.
func (belt *PartyBelt) RemoveListener(listener PartyBeltUpdateListener) {
	for i, l := range belt.listeners {
		if l == listener {
			belt.listeners = append(belt.listeners[:i], belt.listeners[i+1:]...)
			break
		}
	}
}

// Updated sends a visual update to the Session.
func (listener *PartyBeltSessionListener) Updated(slot int, monster Monster) {
	listener.session.QueueEvent(&transport.SetPartySlot{
		Slot:            byte(slot),
		MonsterID:       uint16(monster.ID),
		Gender:          byte(monster.Gender),
		Coloration:      byte(monster.Coloration),
		StatusCondition: byte(monster.StatusCondition),
	})
}
