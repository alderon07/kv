package store

import (
	"kv/pkg/utils"
	"time"
)

type KV[K comparable, V any] struct {
	GMap utils.GMap[K, V]
}

func New[K comparable, V any](defaultTTL time.Duration, cleanUpInterval time.Duration) *KV[K, V]{
	if defaultTTL <= 0 {
		defaultTTL = utils.DEFAULT_TTL
	}

	if cleanUpInterval <= 0 {
		cleanUpInterval = utils.DEFAULT_CLEANUP_INTERVAL
	}

	return &KV[K, V]{
		GMap: *utils.New[K, V](defaultTTL, cleanUpInterval),
	}
}

func (kv *KV[K, V]) Set(key K, value V){
	kv.GMap.Set(key, value)
}

func (kv *KV[K, V]) SetWithTTL(key K, value V, expiresIn time.Duration){
	kv.GMap.SetWithTTL(key, value, expiresIn)
}

func (kv *KV[K, V]) Get(key K) (V, bool){
	item, ok := kv.GMap.Get(key)
	return item.Value, ok
}

func (kv *KV[K, V]) Delete(key K){
	kv.GMap.Delete(key)
}

func (kv *KV[K, V]) SetMultiple(keys map[K]utils.GMapItem[V]){
	kv.GMap.SetMultiple(keys)
}

func (kv *KV[K, V]) GetMultiple(keys []K) map[K]utils.GMapItem[V]{
	return kv.GMap.GetMultiple(keys)
}