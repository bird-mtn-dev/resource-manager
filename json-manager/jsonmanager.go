package jsonmanager

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"

	"golang.org/x/exp/maps"
)

type JSONManager[T any] struct {
	FS        fs.FS
	jsonCache map[string]T
}

func Create[T any]() *JSONManager[T] {
	return CreateWithFS[T](os.DirFS("."))
}

func CreateWithFS[T any](filesystem fs.FS) *JSONManager[T] {
	return &JSONManager[T]{FS: filesystem, jsonCache: make(map[string]T)}
}

func (jm *JSONManager[T]) GetJSON(path string) (T, error) {
	result, exists := jm.jsonCache[path]
	if exists {
		return result, nil
	}
	file, err := jm.FS.Open(path)
	if err != nil {

		return result, err
	}
	dat, err := io.ReadAll(file)
	if err != nil {
		return result, err
	}
	return jm.GetJSONBytes(path, dat)
}

func (jm *JSONManager[T]) GetJSONBytes(key string, dat []byte) (T, error) {
	result, exists := jm.jsonCache[key]
	if exists {
		return result, nil
	}

	err := json.Unmarshal(dat, result)
	if err != nil {
		return result, err
	}
	jm.jsonCache[key] = result
	return result, nil
}

func (fm *JSONManager[T]) Put(key string, value T) {
	fm.jsonCache[key] = value
}

func (fm *JSONManager[T]) Get(key string) T {
	return fm.jsonCache[key]
}

func (fm *JSONManager[T]) Remove(key string) {
	delete(fm.jsonCache, key)
}

// This function will remove all json data in this Manager.
func (fm *JSONManager[T]) Clear() {
	maps.Clear(fm.jsonCache)
}
