package ecs

import (
	"time"
)

// emptyEntityList is an empty slice of Entity's.
var emptyEntityList []*Entity

// systemListener listens for systems being added and removed
// from the world.
type systemListener interface {
	systemAdded(system *System)
	systemRemoved(system *System)
}

// System processes game logic for a set of specific components
// of entities that are indirectly subscribed to the System.
type System struct {
	Policy    SystemPolicy
	Processor Processor

	entities []*Entity
}

// Processor processes all of the entities.
type Processor interface {
	AddedToWorld(world *World) error
	RemovedFromWorld(world *World) error

	Update(world *World, deltaTime time.Duration) error
	Components() ComponentTag
}

// SystemPolicy is a policy of a System to decide if it is the
// right time for the System to update itself and its entities.
type SystemPolicy interface {
	Update(deltaTime time.Duration) bool
}

// NoPolicy is a type of SystemPolicy that says that the System
// is to always run.
type NoPolicy struct{}

// IntervalPolicy is a type of SystemPolicy that accumulates time which
// can be used to decide whether it is time to run a system.
type IntervalPolicy struct {
	rate        time.Duration
	accumulator time.Duration
}

// systemAddition is the addition of a system to the world.
type systemAddition struct {
	system *System
}

// systemRemoval is the removal of a system from the world.
type systemRemoval struct {
	system *System
}

// systemManager manages all of the systems within our world.
type systemManager struct {
	systems         []*System
	systemsToAdd    []systemAddition
	systemsToRemove []systemRemoval
	listeners       []systemListener
}

// NewSystem constructs a new System that is to be updated according
// to the specified policy and when an update is called, the given
// Processor is to take care of that update.
func NewSystem(policy SystemPolicy, processor Processor) *System {
	return &System{
		Policy:    policy,
		Processor: processor,
	}
}

// NewIntervalPolicy creates a new instance of IntervalPolicy.
func NewIntervalPolicy(rate time.Duration) *IntervalPolicy {
	return &IntervalPolicy{rate: rate}
}

// NewDefaultSystemPolicy returns a SystemPolicy that indicates that a
// System is to always run when called for.
func NewDefaultSystemPolicy() SystemPolicy {
	return &NoPolicy{}
}

// newSystemManager constructs a new instance of a manager of entity systems.
func newSystemManager() *systemManager {
	return &systemManager{}
}

// Update returns whether the System should run according to the
// SystemPolicy, which is only true if the amount of accumulated
// time exceeds that of the interval rate.
func (interval *IntervalPolicy) Update(deltaTime time.Duration) bool {
	interval.accumulator += deltaTime

	timeToRun := interval.accumulator >= interval.rate
	if timeToRun {
		interval.accumulator = 0
	}

	return timeToRun
}

// Update returns whether the System should run according to the
// SystemPolicy, which is always true for the NoPolicy implementation.
func (policy *NoPolicy) Update(deltaTime time.Duration) bool {
	// nothing

	return true
}

// add schedules the given system to be added to this system manager.
func (manager *systemManager) add(system *System) {
	manager.systemsToAdd = append(manager.systemsToAdd, systemAddition{system: system})
}

// remove schedules the given system to be removed from this system manager.
func (manager *systemManager) remove(system *System) {
	manager.systemsToRemove = append(manager.systemsToRemove, systemRemoval{system: system})
}

// update updates every registered system to the given World.
func (manager *systemManager) update(world *World, deltaTime time.Duration) error {
	for _, removal := range manager.systemsToRemove {
	Removal:
		for i, system := range manager.systems {
			if system == removal.system {
				manager.systems = append(manager.systems[:i], manager.systems[i+1:]...)
				manager.notifySystemRemoved(removal.system)

				if err := system.Processor.RemovedFromWorld(world); err != nil {
					return err
				}

				break Removal
			}
		}

		manager.systemsToRemove = manager.systemsToRemove[1:]
	}

	for _, addition := range manager.systemsToAdd {
		manager.systems = append(manager.systems, addition.system)
		manager.notifySystemAdded(addition.system)

		manager.systemsToAdd = manager.systemsToAdd[1:]

		if err := addition.system.Processor.AddedToWorld(world); err != nil {
			return err
		}
	}

	for _, system := range manager.systems {
		if system.Policy.Update(deltaTime) {
			if err := system.Processor.Update(world, deltaTime); err != nil {
				return err
			}
		}
	}

	return nil
}

// notifySystemAdded notifies all listener's that the given System
// having been added to this system manager.
func (manager *systemManager) notifySystemAdded(system *System) {
	for _, listener := range manager.listeners {
		listener.systemAdded(system)
	}
}

// notifySystemRemoved notifies all listener's that the given System
// having been removed from this system manager.
func (manager *systemManager) notifySystemRemoved(system *System) {
	for _, listener := range manager.listeners {
		listener.systemRemoved(system)
	}
}

// install installs the given listener into this system manager.
func (manager *systemManager) install(listener systemListener) {
	manager.listeners = append(manager.listeners, listener)
}

// uninstall uninstalls the given listener from this system manager.
func (manager *systemManager) uninstall(listener systemListener) {
	for i, l := range manager.listeners {
		if l == listener {
			manager.listeners = append(manager.listeners[:i], manager.listeners[i+1:]...)
			break
		}
	}
}
