package gmap

import (
	"sync"
	// "errors"
	"context"
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
    CleanUpInterval time.Duration
    Ctx context.Context
    Cancel context.CancelFunc
}

func New[K comparable, V any](defaultTTL time.Duration, cleanUpInterval time.Duration) *GMap[K, V]{
  ctx, cancel := context.WithCancel(context.Background())
  var gmap *GMap[K, V] = &GMap[K, V]{
                          GMap:            make(map[K]GMapItem[V]),
                          DefaultTTL:      defaultTTL,
                          CleanUpInterval: cleanUpInterval,
                          Ctx: ctx,
                          Cancel: cancel,
                        }
  

  go gmap.startCleanUp()

  return gmap
}

func (myGMap *GMap[K, V]) startCleanUp(){
  for{
    myGMap.cleanStaleItems()

    select {
      // wait for cleanUp amount of time
      // When triggered: the case body runs (which is empty), then the loop continues
      case <-time.After(myGMap.CleanUpInterval):
      case <-myGMap.Ctx.Done():
        return
    }
  }
}

func (myGMap *GMap[K, V]) findStaleItems() []K {
  myGMap.Mutex.RLock()
  defer myGMap.Mutex.RUnlock()

  staleKeys := []K{}
  for key , item := range myGMap.GMap{
    if time.Now().After(item.ExpiresAt) {
      staleKeys = append(staleKeys, key)
    }
  }

  return staleKeys
}

func (myGMap * GMap[K, V]) cleanStaleItems(){
  staleKeysToRemove := myGMap.findStaleItems()
  
  myGMap.Mutex.Lock()
  defer myGMap.Mutex.Unlock()
  
  for _, key := range staleKeysToRemove {
    delete(myGMap.GMap, key)
  }
}

// shutdown method
// Cancels the context. Signals the cleanup goroutine to stop. Allows graceful termination.
func (myGMap *GMap[K, V]) Close() {
  myGMap.Cancel()  // ← Calls the cancel function
}

// CRUD
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