package ecs

import "testing"

type sampleComponent struct{}

func (sample *sampleComponent) Tag() ComponentTag {
	return 1 << 0
}

func TestComponentStorage_Add(t *testing.T) {
	storage := newComponentStorage()

	id := newEntityId()
	component := &sampleComponent{}

	storage.add(id, component)
	if c := storage.getComponent(id, component.Tag()); c == nil {
		t.Error("expected component to not be null")
	}
}

func TestComponentStorage_Remove(t *testing.T) {
	storage := newComponentStorage()

	id := newEntityId()
	component := &sampleComponent{}

	storage.add(id, component)
	if c := storage.getComponent(id, component.Tag()); c == nil {
		t.Error("expected component to not be null")
	}

	storage.remove(id, component.Tag())
	if c := storage.getComponent(id, component.Tag()); c != nil {
		t.Error("expected component to be null")
	}
}
