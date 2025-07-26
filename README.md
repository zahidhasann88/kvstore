# Simple In-Memory Key-Value Store

A Redis-like key-value store implementation in Go with CLI interface and TCP server support.

## Features

- **Basic Operations**: SET, GET, DEL commands
- **TTL Support**: Set expiration time for keys using EX parameter
- **Persistence**: Save/Load data to/from JSON files
- **CLI Interface**: Interactive command-line interface
- **TCP Server**: Network interface for remote connections
- **Memory Management**: Automatic cleanup of expired keys
- **Thread Safe**: Concurrent access support with mutex locks
- **Quote Handling**: Proper parsing of quoted values with spaces

## Installation

### Method 1: Direct Installation (Recommended)
```bash
go install github.com/zahidhasann88/kvstore@latest
```

### Method 2: From Source
```bash
git clone https://github.com/zahidhasann88/kvstore.git
cd kvstore
go mod tidy
go build .
```

### Method 3: Run Without Installing
```bash
git clone https://github.com/zahidhasann88/kvstore.git
cd kvstore
go run .
```

## Quick Start

### CLI Mode (Default)
After installation, simply run `kvstore` or `./kvstore` if built from source.

```
> SET name "John Doe"
OK
> GET name
"John Doe"
> SET temp "expires soon" EX 10
OK
> DEL name
1
> SAVE backup.json
OK
> EXIT
```

### TCP Server Mode

**Start server:**
```bash
kvstore server
```

**Connect client:**
```bash
# Method 1: Built-in client
kvstore client

# Method 2: Using telnet
telnet localhost 8080
```

**Example session:**
```bash
# Terminal 1 - Server
$ kvstore server
Starting KV Store TCP Server on :8080...

# Terminal 2 - Client
$ kvstore client
Connected to server at localhost:8080
> SET user "Alice"
< OK
> GET user
< "Alice"
> EXIT
```

## Commands Reference

| Command | Syntax | Description | Example |
|---------|--------|-------------|---------|
| SET | `SET key value [EX seconds]` | Set key-value with optional TTL | `SET user "Alice" EX 30` |
| GET | `GET key` | Retrieve value by key | `GET user` |
| DEL | `DEL key` | Delete a key | `DEL user` |
| SAVE | `SAVE filename` | Export data to JSON file | `SAVE data.json` |
| LOAD | `LOAD filename` | Import data from JSON file | `LOAD data.json` |
| EXIT | `EXIT` or `QUIT` | Close the application | `EXIT` |

## Usage as Go Library

```go
package main

import (
    "fmt"
    "time"
    "github.com/zahidhasann88/kvstore/store"
)

func main() {
    kvStore := store.NewStore()
    defer kvStore.Close()
    
    // Set values
    kvStore.Set("user", "Alice", 0) // No expiration
    kvStore.Set("session", "abc123", 30*time.Second) // 30 second TTL
    
    // Get values
    if value, exists := kvStore.Get("user"); exists {
        fmt.Printf("User: %s\n", value)
    }
    
    // Save to file
    err := kvStore.SaveToFile("data.json")
    if err != nil {
        fmt.Printf("Save error: %v\n", err)
    }
}
```

## Use Cases & Examples

### Development Cache
```bash
kvstore server &  # Background server for development
```

### Session Storage with TTL
```bash
> SET session:user123 "active" EX 3600  # 1 hour session
> SET session:user456 "active" EX 1800  # 30 min session
> SAVE sessions.json  # Backup sessions
```

### Configuration Management
```bash
> SET app:debug "true"
> SET app:max_connections "100"
> SET app:api_key "secret123"
> SAVE config.json
```

### Testing with JSON Data
```bash
> SET test_user1 '{"name":"Alice","email":"alice@test.com"}'
> SET test_user2 '{"name":"Bob","email":"bob@test.com"}'
> GET test_user1
```

### TTL Demonstration
```bash
> SET session "abc123" EX 30
OK
> GET session
"abc123"
# Wait 30 seconds...
> GET session
(nil)
```

### Persistence Workflow
```bash
> SET user1 "Alice"
> SET user2 "Bob"
> SAVE users.json
> DEL user1
> LOAD users.json    # Restores user1
> GET user1
"Alice"
```

## Architecture

### Project Structure
```
kvstore/
├── main.go              # Entry point and CLI logic
├── parser/              # Command parsing
│   └── command.go
├── store/               # Core storage engine
│   ├── store.go         # Main store implementation
│   ├── expiration.go    # TTL management
│   └── persistence.go   # File I/O operations
├── server/              # Network interface
│   └── tcp.go           # TCP server/client
├── utils/               # Helper functions
│   └── helpers.go
└── Makefile            # Build automation
```

### Key Components
- **Store**: Thread-safe in-memory storage with automatic expiration
- **Parser**: Robust command parsing with quote handling
- **Server**: Multi-client TCP server with connection management
- **Persistence**: JSON-based data serialization
- **Expiration**: Timer-based automatic key cleanup

### Technical Details

**Data Structure:**
```go
type Item struct {
    Value     string    `json:"value"`
    ExpiresAt time.Time `json:"expires_at"`
    HasTTL    bool      `json:"has_ttl"`
}
```

**Performance Characteristics:**
- Memory: O(n) where n is number of stored keys
- Operations: O(1) average time complexity for SET/GET/DEL
- Concurrency: Read-write mutex for thread safety
- Persistence: On-demand file operations

**Network Protocol:**
- Plain text commands over TCP
- Line-based protocol (commands end with \n)
- Multi-client support with goroutines

## Building & Development

```bash
make build    # Build binary
make run      # Run CLI mode
make server   # Run TCP server
make test     # Run tests
make clean    # Clean build artifacts
make install  # Install to GOPATH/bin
```

## Configuration

**Default Settings:**
- TCP Server Port: `:8080`
- Max Key Length: 250 characters
- File Format: JSON
- Default TTL: No expiration

## Docker Usage

```dockerfile
FROM golang:1.22-alpine
RUN go install github.com/zahidhasann88/kvstore@latest
EXPOSE 8080
CMD ["kvstore", "server"]
```

```bash
docker build -t kvstore .
docker run -p 8080:8080 kvstore
```

## Troubleshooting

### Installation Issues
```bash
# Check Go version (should be 1.21+)
go version

# Clean and reinstall
go clean -modcache
go install github.com/zahidhasann88/kvstore@latest
```

### Server Issues
```bash
# Check if port is in use
netstat -an | grep :8080

# Kill existing server
pkill kvstore
```

### Permission Issues
```bash
# Make binary executable
chmod +x kvstore

# Ensure GOPATH/bin is in PATH
echo $PATH | grep $(go env GOPATH)/bin
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.