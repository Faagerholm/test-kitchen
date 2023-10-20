package store

import (
	"errors"

	"github.com/faagerholm/page/session"
)

// in-memory store -> connect to reddis or something
var (
	local  map[string]int
	global int
)
var ErrStoreNotInitialized = errors.New("Store not initialized")

func NewCounter() {
	local = make(map[string]int)
}

func Get(id string) (g, s int) {
	if c, ok := local[id]; ok {
		return global, c
	}
	return global, 0
}

func IncrementGlobal() {
	global += 1
	session.Boardcast(global)
}

func IncrementSession(id string) (int, error) {
	if local == nil {
		return 0, ErrStoreNotInitialized
	}

	if c, ok := local[id]; ok {
		c += 1
		local[id] = c
		return c, nil
	}
	local[id] = 1
	return 1, nil
}
