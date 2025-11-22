package main

import (
	"flag"
	"fmt"
	"bufio"
	"os"
	"strings"
	"os"
	"time"
	"os/signal"
	"syscall"

	"kv/internal/store"
)

func main(){
	var defaultTTL 		time.Duration = store.DEFAULT_TTL
	var defaultCleanup 	time.Duration = store.DEFAULT_CLEANUP_INTERVAL

	ttl 		:= flag.Duration("ttl", defaultTTL , "default TTL for keys")
	cleanup 	:= flag.Duration("cleanup", defaultCleanup , "default cleanup time for expired keys")
	// port 		:= flag.Int("port", 8080, "port number")
	// host 		:= flag.String("host", "127.0.0.1", "host")
	
	flag.Usage = usage
	flag.Parse()
	
	fmt.Println("TTL:", *ttl)
	fmt.Println("Cleanup:", *cleanup)

	kv := store.NewWithConfig[string, string](*ttl, *cleanup)
	fmt.Printf("KV Store created: %v\n", kv)
	defer kv.Close()
	fmt.Println("KV Store running... Press Ctrl+C to exit")
    
	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	//signal.Notify() - Tells Go's signal package to forward signals to your channel
	// os.Interrupt - Catches Ctrl+C (SIGINT on Unix, equivalent on Windows)
	// syscall.SIGTERM - Catches termination signal (commonly sent by kill command or container orchestrators
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	// Blocks execution until a signal is received. This is a blocking receive operation on the channel
	<-sigChan
	
	fmt.Println("\nShutting down gracefully...")
}

func usage() {
	fmt.Println(`kv - An in-memory key-value store

		Usage:
			kv [flags]

		Examples:
			kv -ttl=5m -cleanup=30s
			kv -port=6379 -host=0.0.0.0
			kv -help

		Flags:
			-ttl      default time to live for keys in the kv store (default: 10m)
			-cleanup  default cleanup interval to delete keys with expired ttl (default: 1h)
			-port     port to run the server on (default: 8080)
			-host     host on which the server would run (default: 127.0.0.1)
			-help     show this usage information

		Duration Formats:
			"5s"    - 5 seconds
			"10m"   - 10 minutes  
			"2h"    - 2 hours
			"1h30m" - 1 hour 30 minutes
			"100ms" - 100 milliseconds`)
	fmt.Println("STORE CREATED")
	kv := store.New[string, string]()

	startREPL(kv)
}

func startREPL(store *store.KV[string,string]){
	scanner := bufio.NewScanner(os.Stdin)

	for{
		fmt.Print("> ")

		if !scanner.Scan() {
			break // EOF or error
		}

		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		cmd := strings.ToUpper(parts[0])
		var key string
		var value string

		switch cmd {
			case "GET":
				key = parts[1]
				val, exists := store.Get(key)
				if !exists {
					fmt.Printf("Key %s Does Not Exist\n", key)
					break
				}
				fmt.Println(val)
			case "SET":
				key = parts[1]
				value = parts[2]
				store.Set(key, value)
				fmt.Printf("Set Value %s for Key %s\n", value, key)
			case "DELETE":
				key = parts[1]
				store.Delete(key)
				fmt.Printf("Delted Key %s\n", key)
			case "LIST":
				fmt.Println(store.List())
			case "EXIT", "QUIT":
				return
		}
	}
}