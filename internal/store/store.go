package store

import (
	"kv/pkg/gmap"
	"time"
)

const DEFAULT_TTL 							= gmap.DEFAULT_TTL
const DEFAULT_CLEANUP_INTERVAL 	= gmap.DEFAULT_CLEANUP_INTERVAL

type KV[K comparable, V any] struct {
  GMap gmap.GMap[K, V]
}

// New creates a KV store with default TTL and cleanup interval.
func New[K comparable, V any]() *KV[K, V]{
  return NewWithConfig[K, V](DEFAULT_TTL, DEFAULT_CLEANUP_INTERVAL)
}

// NewWithConfig creates a KV store with the provided TTL and cleanup interval.
// Zero or negative values fall back to the underlying gmap defaults.
func NewWithConfig[K comparable, V any](defaultTTL time.Duration, cleanUpInterval time.Duration) *KV[K, V]{
  return &KV[K, V]{
    GMap: *gmap.NewWithConfig[K, V](defaultTTL, cleanUpInterval),
  }
}

func (kv *KV[K, V]) Close(){
  kv.GMap.Close()
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

func (kv *KV[K, V]) SetMultiple(keys map[K]gmap.GMapItem[V]){
  kv.GMap.SetMultiple(keys)
}

func (kv *KV[K, V]) GetMultiple(keys []K) map[K]gmap.GMapItem[V]{
  return kv.GMap.GetMultiple(keys)
}

func (kv *KV[K, V]) Size() int {
  return kv.GMap.Size()
}

func (kv *KV[K, V]) Clear() {
  kv.GMap.Clear()
}

func (kv *KV[K, V]) List() string {
  return kv.GMap.List()
}