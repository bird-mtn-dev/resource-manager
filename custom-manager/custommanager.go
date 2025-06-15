package custommanager

import (
	"golang.org/x/exp/maps"
)

type CustomManager[T any] struct {
	customCache map[string]T
}

func Create[T any]() *CustomManager[T] {
	return &CustomManager[T]{customCache: make(map[string]T)}
}

func (fm *CustomManager[T]) Put(key string, value T) {
	fm.customCache[key] = value
}

func (fm *CustomManager[T]) Get(key string) T {
	return fm.customCache[key]
}

func (fm *CustomManager[T]) Remove(key string) {
	delete(fm.customCache, key)
}

// This function will remove all custom data in this Manager
func (fm *CustomManager[T]) Clear() {
	maps.Clear(fm.customCache)
}
