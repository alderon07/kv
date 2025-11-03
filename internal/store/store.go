package store

import (
	"kv/pkg/utils"
	"time"
)

type KVStore[K comparable, V any] struct {
	gMap utils.GMap[K, V]
	createdAt time.Time
	// ttl time.Duration
}

func (store *KVStore[K, V]) NewKVStore(){
	store.createdAt = time.Now()
	store.gMap = *utils.NewGMap[K,V]()
}

func (store *KVStore[K, V]) Set(key K, value V){
	store.gMap.Set(key, value)
}

func (store *KVStore[K, V]) Get(key K) (V, bool){
	return store.gMap.Get(key)
}

func (store *KVStore[K, V]) Delete(key K){
	store.gMap.Delete(key)
}

func (store *KVStore[K, V]) SetMultiple(keys map[K]V){
	store.gMap.SetMultiple(keys)
}

func (store *KVStore[K, V]) GetMultiple(keys []K) map[K]V{
	return store.gMap.GetMultiple(keys)
}