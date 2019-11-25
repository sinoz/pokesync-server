package ecs

// System processes game logic for a set of specific components
// of entities that are indirectly subscribed to the System.
type System struct {
	Policy    SystemPolicy
	Processor Processor
}

// Processor processes all of the entities.
type Processor interface {
	Update(entities []*Entity, deltaTime float64)
}

// SystemPolicy is a policy of a System to decide if it is the
// right time for the System to update itself and its entities.
type SystemPolicy interface {
	ShouldRun() bool
}

// NoPolicy is a type of SystemPolicy that says that the System
// is to always run.
type NoPolicy struct{}

// systemAddition is the addition of a system to the engine.
type systemAddition struct {
	system System
}

// systemRemoval is the removal of a system from the engine.
type systemRemoval struct {
	system System
}

// systemManager manages all of the systems within our engine.
type systemManager struct {
	systems         []System
	systemsToAdd    []systemAddition
	systemsToRemove []systemRemoval
}

// NoSystemPolicy returns a SystemPolicy that indicates that a
// System is to always run when called for.
func NoSystemPolicy() SystemPolicy {
	return &NoPolicy{}
}

// newSystemManager constructs a new instance of a manager of entity systems.
func newSystemManager() *systemManager {
	return &systemManager{}
}

// ShouldRun returns whether the System should run according to the
// SystemPolicy, which is always true for the NoPolicy implementation.
func (policy *NoPolicy) ShouldRun() bool {
	return true
}

// add schedules the given system to be added to this system manager.
func (manager *systemManager) add(system System) {
	manager.systemsToAdd = append(manager.systemsToAdd, systemAddition{system: system})
}

// remove schedules the given system to be removed from this system manager.
func (manager *systemManager) remove(system System) {
	manager.systemsToRemove = append(manager.systemsToRemove, systemRemoval{system: system})
}

// update updates every registered system.
func (manager *systemManager) update(deltaTime float64) {
	for _, removal := range manager.systemsToRemove {
	Removal:
		for i, system := range manager.systems {
			if system == removal.system {
				manager.systems = append(manager.systems[:i], manager.systems[i+1:]...)
				break Removal
			}
		}

		manager.systemsToRemove = manager.systemsToRemove[1:]
	}

	for _, addition := range manager.systemsToAdd {
		manager.systems = append(manager.systems, addition.system)
		manager.systemsToAdd = manager.systemsToAdd[1:]
	}

	for _, system := range manager.systems {
		if system.Policy.ShouldRun() {
			// TODO
		}
	}
}
