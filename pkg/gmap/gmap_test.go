package gmap

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestSet(t *testing.T){
	kv := New[string, string]()
	defer kv.Close()

	// add one k, v pair
	kv.Set("1", "You")
	item, ok := kv.Get("1")
	if !ok {
		t.Errorf("TestGet KEY_NOT_FOUND")
	}
	result := item.Value
	expected := "You"
	if result != expected {
		t.Errorf("kv.Set() returned %s expected %s", result, expected)
	}

	// add another
	kv.Set("2", "Suck")
	item, ok = kv.Get("2")
	if !ok {
		t.Errorf("TestGet KEY_NOT_FOUND")
	}
	result = item.Value
	expected = "Suck"
	if result != expected {
		t.Errorf("kv.Set() returned %s expected %s", result, expected)
	}

	// update the first key
	kv.Set("1", "I")
	item, ok = kv.Get("1")
	if !ok {
		t.Errorf("TestGet KEY_NOT_FOUND")
	}
	result = item.Value
	expected = "I"
	if result != expected {
		t.Errorf("kv.Set() returned %s expected %s", result, expected)
	}
}

func TestGet(t *testing.T){
	kv := New[string, string]()
	defer kv.Close()

	// add one k, v pair
	kv.Set("1", "You")
	item, ok := kv.Get("1")
	if !ok {
		t.Errorf("TestGet KEY_NOT_FOUND")
	}
	result := item.Value
	expected := "You"
	if result != expected {
		t.Errorf("TestGet expected %s result %s", expected, result)
	}
}

func TestDelete(t *testing.T){
	kv := New[string, string]()
	defer kv.Close()

	// add one k, v pair
	kv.Set("1", "You")
	kv.Set("2", "Suck")
	kv.Delete("1")
	_ , ok := kv.Get("1")
	if ok {
		t.Errorf("TestDelete KEY_FOUND_AFTER_DELETE")
	}

	item, ok := kv.Get("2")
	if !ok {
		t.Errorf("TestDelete KEY_NOT_FOUND")
	}
	result := item.Value
	expected := "Suck"
	if result != expected {
		t.Errorf("kv.Set() returned %s expected %s", result, expected)
	}
}

func TestClear(t *testing.T){
	kv := New[string, string]()
	defer kv.Close()

	kv.Set("1", "You")
	kv.Set("2", "Suck")
	kv.Delete("1")
	kv.Clear()

	expected := 0
	result := kv.Size()
	if expected != result {
		t.Errorf("TestClear CLEAR_DID_NOT_DELETE_ALL_ITEMS")
	}
}

func TestSize(t *testing.T){
	kv := New[string, string]()
	defer kv.Close()

	kv.Set("1", "You")
	kv.Set("2", "Suck")
	expected := 2
	result := kv.Size()
	if expected != result {
		t.Errorf("TestSize WRONG_SIZE expected: %d result %d", expected, result)
	}
}

func TestSetWithTTL(t *testing.T){
	kv := New[string, string]()
	defer kv.Close()

	kv.SetWithTTL("1", "You", 2 * time.Second)
	
	item, _ := kv.Get("1")
	result := item.Value
	expected := "You"

	if expected != result {
		t.Errorf("TestSetWithTTL expected: %s result %s", expected, result)
	}

	time.Sleep(2 * time.Second)

	_, ok := kv.Get("1")
	if !ok {
		t.Errorf("TestSetWithTTL ITEM_EXISTS_AFTER_TTL_EXPIRATION")
	}
}

func TestConcurrentSets(t *testing.T){
	kv := New[int, string]()
	
	var wg sync.WaitGroup
	count := 100
	value := "Value:"
	
	// spawn go routines to run sets concurrently
	for i := range count {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			kv.Set(idx, value + strconv.Itoa(idx))
		}(i)
	}

	wg.Wait()

	// check if the keys exists
	for i := range count {
		item, ok := kv.Get(i)
		result := item.Value
		if !ok {
			t.Errorf("TestConcurrentSets key %d does not exist", i)
		}
		
		expected := value + strconv.Itoa(i)
		if expected != result {
			t.Errorf("TestConcurrentSets value %s did not match expected %s for key %d", result, expected, i)
		}
	} 
}

func TestConcurrentGets(t *testing.T){
	kv := New[int, string]()
	
	count := 100
	value := "Value:"
	// set key/values
	for i := range count {
		kv.Set(i, value + strconv.Itoa(i))
	}
	
	var wg sync.WaitGroup
	numOfGoroutines := 500
	// spawn go routines to run gets concurrently
	for i := range numOfGoroutines {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			
			key := idx % count
			
			item, ok := kv.Get(key)
			result := item.Value
			if !ok {
				t.Errorf("key %d does not exist", key)
			}

			expected := value + strconv.Itoa(key)
			if expected != result {
				t.Errorf("TestConcurrentSets value %s did not match expected %s for key %d", result, expected, key)
			}

		}(i)
	}

	wg.Wait()
}

func TestConcurrentSetGetDelete(t *testing.T){
	return
}