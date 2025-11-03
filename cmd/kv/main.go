package main

import (
	"fmt"
	"kv/internal/store"
)

func main(){
	fmt.Println("kv")
	store := make(store.KVStore)
}