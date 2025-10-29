package utils

import(
	"sync"
	// "errors"
)

type GMap[K comparable, V any] struct {
	mutex sync.RWMutex
	gMap map[K]V 
}

func New[K comparable, V any]() *GMap[K, V]{
	return &GMap[K, V]{gMap: make(map[K]V)}
}

func (myGMap *GMap[K, V]) Set(key K, value V) {
	myGMap.mutex.Lock()
	defer myGMap.mutex.Unlock()

	myGMap.gMap[key] = value
}

func (myGMap *GMap[K, V]) Delete(key K) {
	myGMap.mutex.Lock()
	defer myGMap.mutex.Unlock()
	delete(myGMap.gMap, key)
}

func (myGMap *GMap[K, V]) Get(key K) (V , bool) {
	myGMap.mutex.RLock()
	defer myGMap.mutex.RUnlock()
	v, ok := myGMap.gMap[key]

	return v, ok
}

func (myGMap *GMap[K, V]) Clear() {
	myGMap.mutex.Lock()
	defer myGMap.mutex.Unlock()

	myGMap.gMap = make(map[K]V)
}

func (myGMap *GMap[K, V]) Length() int {
	myGMap.mutex.RLock()
	defer myGMap.mutex.RUnlock()

	return len(myGMap.gMap)
}

func (myGMap *GMap[K, V]) GetMultiple(keys []K) map[K]V{
	myGMap.mutex.RLock()
	defer myGMap.mutex.RUnlock()

	resultMap := make(map[K]V)
	for _, key := range keys {
		v, ok := myGMap.gMap[key]
		if(ok){
			resultMap[key] = v
		}
	}
	return resultMap
}

func (myGMap *GMap[K, V]) SetMultiple(){
	
}