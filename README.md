# KV - In-Memory Key-Value Store

A high-performance, thread-safe, in-memory key-value store written in Go with automatic TTL expiration and cleanup.

## Features

- ✨ **Generic Type Support** - Use any comparable type as key and any type as value
- 🔒 **Thread-Safe** - Built-in concurrency support with RWMutex
- ⏰ **Automatic TTL Expiration** - Keys automatically expire after a configurable time
- 🧹 **Background Cleanup** - Periodic garbage collection of expired keys
- 🚀 **High Performance** - Optimized for concurrent read/write operations
- 🛡️ **Graceful Shutdown** - Clean resource cleanup on termination
- 🎯 **Batch Operations** - Support for multiple get/set operations

## Architecture

The project follows a clean architecture with three main layers:

```
kv/
├── cmd/kv/              # Application entry point
│   └── main.go          # CLI and signal handling
├── internal/store/      # Store abstraction layer
│   └── store.go         # Public-facing KV store API
└── pkg/gmap/            # Core generic map implementation
    ├── gmap.go          # Thread-safe generic map with TTL
    └── gmap_test.go     # Comprehensive test suite
```

### Components

- **`cmd/kv`**: Command-line interface with flag parsing and graceful shutdown
- **`internal/store`**: High-level KV store API that wraps the generic map
- **`pkg/gmap`**: Core implementation with TTL management and automatic cleanup

## Installation

### Prerequisites

- Go 1.22 or higher (uses modern Go features like range over int)

### Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd kv

# Build the binary
make build

# Or run directly
make run
```

## Usage

### Command Line

Run the KV store with default settings:

```bash
./bin/kv
```

Run with custom TTL and cleanup interval:

```bash
./bin/kv -ttl=5m -cleanup=30s
```

Show help:

```bash
./bin/kv -help
```

### Available Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-ttl` | Default time to live for keys | `1h` |
| `-cleanup` | Cleanup interval for expired keys | `24h` |
| `-help` | Show usage information | - |

### Duration Formats

Go's standard duration formats are supported:

- `5s` - 5 seconds
- `10m` - 10 minutes
- `2h` - 2 hours
- `1h30m` - 1 hour 30 minutes
- `100ms` - 100 milliseconds

### Programmatic Usage

#### Basic Operations

```go
import "kv/pkg/gmap"

// Create a new KV store with defaults (1h TTL, 24h cleanup)
store := gmap.New[string, string]()
defer store.Close()

// Set a value
store.Set("username", "alice")

// Get a value
item, ok := store.Get("username")
if ok {
    fmt.Println(item.Value) // Output: alice
}

// Delete a value
store.Delete("username")

// Clear all values
store.Clear()

// Get size
size := store.Size()
```

#### Custom Configuration

```go
import (
    "kv/pkg/gmap"
    "time"
)

// Create with custom TTL and cleanup interval
store := gmap.NewWithConfig[int, string](
    5 * time.Minute,  // TTL: 5 minutes
    30 * time.Second, // Cleanup: 30 seconds
)
defer store.Close()

// Set with custom TTL
store.SetWithTTL("session", "xyz123", 10*time.Minute)
```

#### Batch Operations

```go
// Get multiple keys at once
keys := []string{"user1", "user2", "user3"}
results := store.GetMultiple(keys)

for key, item := range results {
    fmt.Printf("%s: %s\n", key, item.Value)
}

// Set multiple keys at once
batch := map[string]gmap.GMapItem[string]{
    "key1": {Value: "value1", ExpiresAt: time.Now().Add(1*time.Hour)},
    "key2": {Value: "value2", ExpiresAt: time.Now().Add(2*time.Hour)},
}
store.SetMultiple(batch)
```

#### Generic Types

```go
// String keys, integer values
intStore := gmap.New[string, int]()
intStore.Set("count", 42)

// Integer keys, struct values
type User struct {
    Name string
    Age  int
}
userStore := gmap.New[int, User]()
userStore.Set(1, User{Name: "Alice", Age: 30})

// Custom comparable type as key
type UserID string
customStore := gmap.New[UserID, string]()
customStore.Set(UserID("user123"), "Alice")
```

## API Reference

### GMap Methods

| Method | Description |
|--------|-------------|
| `New[K, V]()` | Create a new GMap with default settings |
| `NewWithConfig[K, V](ttl, cleanup)` | Create with custom TTL and cleanup interval |
| `Set(key, value)` | Set a key-value pair with default TTL |
| `SetWithTTL(key, value, duration)` | Set with custom TTL |
| `Get(key)` | Retrieve a value by key |
| `Delete(key)` | Remove a key-value pair |
| `Clear()` | Remove all key-value pairs |
| `Size()` | Get the number of stored items |
| `GetMultiple(keys)` | Retrieve multiple values at once |
| `SetMultiple(items)` | Set multiple key-value pairs at once |
| `Close()` | Stop background cleanup and cleanup resources |

## Development

### Run Tests

```bash
# Run all tests
make test

# Run tests with verbose output
go test ./... -v

# Run specific test
go test ./pkg/gmap -run TestConcurrentSets -v
```

### Format Code

```bash
make fmt
```

### Vet Code

```bash
make vet
```

### Clean Build Artifacts

```bash
make clean
```

## How It Works

### TTL Management

1. Each key-value pair is stored with an `ExpiresAt` timestamp
2. When you `Set()` a key, it gets a TTL based on the configured default
3. A background goroutine periodically scans for expired keys
4. Expired keys are automatically deleted during cleanup cycles

### Concurrency

- **Read operations** use `RLock()` allowing multiple concurrent readers
- **Write operations** use `Lock()` for exclusive access
- Lock-free operations where possible for maximum throughput
- Tested with hundreds of concurrent goroutines

### Graceful Shutdown

The application handles shutdown signals (`SIGINT`, `SIGTERM`) gracefully:

1. User presses Ctrl+C or system sends SIGTERM
2. Signal is caught by the signal handler
3. `Close()` is called via `defer`
4. Background cleanup goroutine is stopped via context cancellation
5. Resources are cleaned up properly

## Testing

The project includes comprehensive tests covering:

- ✅ Basic CRUD operations
- ✅ TTL expiration behavior
- ✅ Concurrent reads and writes
- ✅ Batch operations
- ✅ Size and clear operations

Run tests with:

```bash
go test ./pkg/gmap -v
```

## Performance Considerations

- **Read-heavy workloads**: Excellent performance due to RWMutex
- **Write-heavy workloads**: Good performance with mutex contention handling
- **Memory**: No built-in memory limits (consider implementing if needed)
- **Cleanup overhead**: Configurable cleanup interval to balance freshness vs. CPU usage

## Future Enhancements

- [ ] HTTP/TCP server interface
- [ ] Persistence to disk
- [ ] Memory limits and eviction policies (LRU, LFU)
- [ ] Metrics and monitoring
- [ ] Clustering and replication
- [ ] Watch/subscribe functionality

## License

[Specify your license here]

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Author

[Your name here]

