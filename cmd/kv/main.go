package main

import (
	"flag"
	"fmt"
	"kv/internal/store"
	"os"
	"time"
	"os/signal"
	"syscall"
)

func main(){
	var defaultTTL 			time.Duration	= 1 * time.Hour
	var defaultCleanup 	time.Duration = 24 * time.Hour
	
	ttl 		:= flag.Duration("ttl", defaultTTL , "default TTL for keys")
	cleanup := flag.Duration("cleanup", defaultCleanup , "default cleanup time for expired keys")
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
}