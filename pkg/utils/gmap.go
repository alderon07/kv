package utils

import (
	"sync"
	// "errors"
	"time"
)

const DEFAULT_TTL time.Duration                 = 24 * time.Hour
const DEFAULT_CLEANUP_INTERVAL time.Duration    = 24 * time.Hour

type GMapItem[V any] struct {
    Value V
    ExpiresAt time.Time
}

type GMap[K comparable, V any] struct {
    Mutex sync.RWMutex
    GMap map[K]GMapItem[V]
    DefaultTTL time.Duration
    cleanupInterval time.Duration
}

func New[K comparable, V any](defaultTTL time.Duration, cleanupInterval time.Duration) *GMap[K, V]{
    return &GMap[K, V]{
        GMap:            make(map[K]GMapItem[V]),
        DefaultTTL:      defaultTTL,
        cleanupInterval: cleanupInterval,
    }
}

func (myGMap *GMap[K, V]) SetWithTTL(key K, value V, expiresIn time.Duration) {
    myGMap.Mutex.Lock()
    defer myGMap.Mutex.Unlock()
    
    myGMap.GMap[key] = GMapItem[V]{
        Value: value,
        ExpiresAt: time.Now().Add(expiresIn),
    }
}

func (myGMap *GMap[K, V]) Set(key K, value V) {
    myGMap.SetWithTTL(key, value, DEFAULT_TTL)
}


func (myGMap *GMap[K, V]) Delete(key K) {
    myGMap.Mutex.Lock()
    defer myGMap.Mutex.Unlock()
    delete(myGMap.GMap, key)
}

func (myGMap *GMap[K, V]) Get(key K) (GMapItem[V], bool) {
    myGMap.Mutex.RLock()
    defer myGMap.Mutex.RUnlock()
    
    item, ok := myGMap.GMap[key]
    // if(!ok || item.ExpiresAt.Before(time.Now())){
    // 	delete(myGMap.GMap, key)
    // 	return GMapItem[V]{}, false
    // }
    return item, ok
}

func (myGMap *GMap[K, V]) Clear() {
    myGMap.Mutex.Lock()
    defer myGMap.Mutex.Unlock()

    myGMap.GMap = make(map[K]GMapItem[V])
}

func (myGMap *GMap[K, V]) Length() int {
    myGMap.Mutex.RLock()
    defer myGMap.Mutex.RUnlock()

    return len(myGMap.GMap)
}

func (myGMap *GMap[K, V]) GetMultiple(keys []K) map[K]GMapItem[V]{
    myGMap.Mutex.RLock()
    defer myGMap.Mutex.RUnlock()

    resultMap := make(map[K]GMapItem[V])
    for _, key := range keys {
        item, ok := myGMap.GMap[key]
        if(ok){
            resultMap[key] = item
        }
    }
    return resultMap
}

func (myGMap *GMap[K, V]) SetMultiple(keys map[K]GMapItem[V])  {
    myGMap.Mutex.Lock()
    defer myGMap.Mutex.Unlock()

    for key, item := range keys {
        myGMap.GMap[key] = GMapItem[V]{
            Value: item.Value,
            ExpiresAt: item.ExpiresAt,
        }
    }
}