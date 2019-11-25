package ecs

// TODO
type Engine struct {
	entityManager    *entityManager
	systemManager    *systemManager
	componentStorage *componentStorage
}

// TODO
func NewEngine() *Engine {
	return &Engine{
		entityManager:    newEntityManager(),
		systemManager:    newSystemManager(),
		componentStorage: newComponentStorage(),
	}
}

// TODO
func (engine *Engine) Update(deltaTime float64) {
	// TODO
}
