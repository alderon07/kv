package main

import (
	"fmt"
	"kv/internal/store"
)

func main(){
	fmt.Println("kv")
	kvStore := store.NewStore[string, string]()
	fmt.Printf("KV Store created: %v\n", kvStore)
}