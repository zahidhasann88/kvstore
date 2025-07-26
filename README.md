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

## Quick Start

### CLI Mode (Default)

```bash
# Clone and run
git clone <https://github.com/zahidhasann88/kvstore.git>
cd kvstore
go run .
```

Example session:
```
> SET name "John Doe"
OK
> GET name
"John Doe"
> SET temp "expires soon" EX 10
OK
> GET temp
"expires soon"
> DEL name
1
> SAVE backup.json
OK
> EXIT
```

### TCP Server Mode

Start server:
```bash
go run . server
```

Connect from another terminal:
```bash
go run . client
# Or use telnet: telnet localhost 8080
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

## Installation & Usage

### From Source
```bash
git clone <your-repo-url>
cd kvstore
go mod tidy
make build    # Creates kvstore binary
./kvstore     # Run CLI mode
./kvstore server  # Run TCP server
```

### As Go Package
```bash
go install github.com/zahidhasann88/kvstore@latest
kvstore       # CLI mode
kvstore server # Server mode
```

### Key Components

- **Store**: Thread-safe in-memory storage with automatic expiration
- **Parser**: Robust command parsing with quote handling
- **Server**: Multi-client TCP server with connection management
- **Persistence**: JSON-based data serialization
- **Expiration**: Timer-based automatic key cleanup

## Building & Development

```bash
make build    # Build binary
make run      # Run CLI mode
make server   # Run TCP server
make test     # Run tests (when added)
make clean    # Clean build artifacts
make install  # Install to GOPATH/bin
```

## Configuration

Default settings:
- TCP Server Port: `:8080`
- Max Key Length: 250 characters
- File Format: JSON
- Default TTL: No expiration

## Performance Characteristics

- **Memory**: O(n) where n is number of stored keys
- **Operations**: O(1) average time complexity for SET/GET/DEL
- **Concurrency**: Read-write mutex for thread safety
- **Persistence**: On-demand file operations

## Examples

### Basic Operations
```bash
> SET counter 1
OK
> SET message "Hello, World!"
OK
> GET counter
"1"
> GET message
"Hello, World!"
```

### TTL (Time To Live)
```bash
> SET session "abc123" EX 30
OK
> GET session
"abc123"
# Wait 30 seconds...
> GET session
(nil)
```

### Persistence
```bash
> SET user1 "Alice"
OK
> SET user2 "Bob"
OK
> SAVE users.json
OK
> DEL user1
1
> LOAD users.json
OK
> GET user1
"Alice"
```

### Network Mode
```bash
# Terminal 1: Start server
$ go run . server
Starting KV Store TCP Server on :8080...
KV Store server listening on :8080

# Terminal 2: Connect client
$ go run . client
Connected to server at localhost:8080
Enter commands (type EXIT to quit):
> SET distributed "works!"
< OK
> GET distributed
< "works!"
```

## Technical Details

### Data Structure
```go
type Item struct {
    Value     string    `json:"value"`
    ExpiresAt time.Time `json:"expires_at"`
    HasTTL    bool      `json:"has_ttl"`
}
```

### Expiration Management
- Timer-based automatic cleanup
- Lazy expiration on access
- Thread-safe timer management

### Network Protocol
- Plain text commands over TCP
- Line-based protocol (commands end with \n)
- Multi-client support with goroutines