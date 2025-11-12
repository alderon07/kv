package main

import (
	"fmt"
	"kv/internal/store"
)

func main(){
	fmt.Println("kv")
	kv := store.New[string, string]
	fmt.Printf("KV Store created: %v\n", kv)
}