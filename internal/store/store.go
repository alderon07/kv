package store

import (
	"kv/pkg/utils"
)


type KVStore[K comparable, V any] struct {
	gMap utils.GMap[K, V]
}