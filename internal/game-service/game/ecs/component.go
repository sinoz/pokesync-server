package ecs

// ComponentTag is a unique bit mask value that is assigned to
// each type of Component for identification purposes.
type ComponentTag int

// Component represents a data structure of a specific domain that
// is to describe a part of an Entity.
type Component interface {
	Tag() ComponentTag
}
