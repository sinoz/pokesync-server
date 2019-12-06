package game

import (
	"errors"
	"gitlab.com/pokesync/game-service/internal/game-service/game/session"
	"gitlab.com/pokesync/game-service/internal/game-service/game/transport"
)

// CoinBagListener listens for changes made to the CoinBag.
type CoinBagListener interface {
	PokeDollarsUpdated(newValue int)
	DonatorPointsUpdated(newValue int)
}

// CoinBagSessionListener is a CoinBagListener that listens
// for changes made to the CoinBag to visually apply these
// changes to the Client as well.
type CoinBagSessionListener struct {
	session *session.Session
}

// CoinBag is a player's bag of coins.
type CoinBag struct {
	Dollars       int
	DonatorPoints int

	listeners []CoinBagListener
}

// NewCoinBag constructs a new instance of a CoinBag.
func NewCoinBag() *CoinBag {
	return &CoinBag{}
}

// AddPokeDollars adds the specified amount of pokedollars to this bag.
// May return an error if the given amount is negative.
func (bag *CoinBag) AddPokeDollars(amount int) error {
	if err := bag.checkBoundaries(amount); err != nil {
		return err
	}

	bag.Dollars += amount
	bag.notifyPokeDollarsUpdated(bag.Dollars)

	return nil
}

// AddDonatorPoints adds the specified amount of donator points to this bag.
// May return an error if the given amount is negative.
func (bag *CoinBag) AddDonatorPoints(amount int) error {
	if err := bag.checkBoundaries(amount); err != nil {
		return err
	}

	bag.DonatorPoints += amount
	bag.notifyDonatorPointsUpdated(bag.DonatorPoints)

	return nil
}

// WithdrawPokeDollars withdraws the specified amount of pokedollars from
// this bag.
func (bag *CoinBag) WithdrawPokeDollars(amount int) {
	bag.Dollars -= amount
	if bag.Dollars < 0 {
		bag.Dollars = 0
	}

	bag.notifyPokeDollarsUpdated(bag.Dollars)
}

// WithdrawDonatorPoints withdraws the specified amount of donator points
// from this bag.
func (bag *CoinBag) WithdrawDonatorPoints(amount int) {
	bag.DonatorPoints -= amount
	if bag.DonatorPoints < 0 {
		bag.DonatorPoints = 0
	}

	bag.notifyDonatorPointsUpdated(bag.DonatorPoints)
}

// notifyPokeDollarsUpdated notifies every subscribed listener
// of the CoinBag's new amount of pokedollars.
func (bag *CoinBag) notifyPokeDollarsUpdated(value int) {
	for _, listener := range bag.listeners {
		listener.PokeDollarsUpdated(value)
	}
}

// notifyPokeDollarsUpdated notifies every subscribed listener
// of the CoinBag's new amount of donator points.
func (bag *CoinBag) notifyDonatorPointsUpdated(value int) {
	for _, listener := range bag.listeners {
		listener.DonatorPointsUpdated(value)
	}
}

// AddListener subscribes the given listener to receive notifications.
func (bag *CoinBag) AddListener(listener CoinBagListener) {
	bag.listeners = append(bag.listeners, listener)
}

// RemoveListener unsubscribes the given listener from receiving notifications.
func (bag *CoinBag) RemoveListener(listener CoinBagListener) {
	for i, l := range bag.listeners {
		if l == listener {
			bag.listeners = append(bag.listeners[:i], bag.listeners[i+1:]...)
			break
		}
	}
}

// checkBoundaries checks if the given amount is valid. If not,
// it returns an error.
func (bag *CoinBag) checkBoundaries(amount int) error {
	if amount < 0 {
		return errors.New("given coin amount is negative")
	}

	return nil
}

// PokeDollarsUpdated sends a visual update to the Session.
func (listener *CoinBagSessionListener) PokeDollarsUpdated(newValue int) {
	listener.session.QueueEvent(&transport.SetPokeDollar{Amount: uint32(newValue)})
}

// DonatorPointsUpdated sends a visual update to the Session.
func (listener *CoinBagSessionListener) DonatorPointsUpdated(newValue int) {
	listener.session.QueueEvent(&transport.SetDonatorPoints{Amount: uint32(newValue)})
}
