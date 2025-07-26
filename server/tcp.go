package server

import (
	"bufio"
	"fmt"
	"github.com/zahidhasann88/kvstore/parser"
	"github.com/zahidhasann88/kvstore/store"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func StartTCPServer(address string) {
	kvStore := store.NewStore()
	defer kvStore.Close()

	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Printf("KV Store server listening on %s\n", address)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		fmt.Println("\nShutting down server...")
		listener.Close()
		os.Exit(0)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		go handleConnection(conn, kvStore)
	}
}

func StartTCPClient(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to server at %s\n", address)
	fmt.Println("Enter commands (type EXIT to quit):")

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Printf("< %s\n", scanner.Text())
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		_, err := conn.Write([]byte(input + "\n"))
		if err != nil {
			fmt.Printf("Failed to send command: %v\n", err)
			break
		}

		if strings.ToUpper(input) == "EXIT" {
			break
		}
	}
}

func handleConnection(conn net.Conn, kvStore *store.Store) {
	defer conn.Close()

	fmt.Printf("Client connected: %s\n", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		cmd, err := parser.ParseCommand(input)
		if err != nil {
			response := fmt.Sprintf("ERROR: %v\n", err)
			conn.Write([]byte(response))
			continue
		}

		result := executeCommand(kvStore, cmd)
		response := result + "\n"
		conn.Write([]byte(response))

		if cmd.Type == parser.EXIT {
			break
		}
	}

	fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr())
}

func executeCommand(kvStore *store.Store, cmd *parser.Command) string {
	switch cmd.Type {
	case parser.SET:
		kvStore.Set(cmd.Key, cmd.Value, cmd.TTL)
		return "OK"
	case parser.GET:
		value, exists := kvStore.Get(cmd.Key)
		if !exists {
			return "(nil)"
		}
		return value
	case parser.DEL:
		deleted := kvStore.Delete(cmd.Key)
		if deleted {
			return "1"
		}
		return "0"
	case parser.SAVE:
		err := kvStore.SaveToFile(cmd.Key)
		if err != nil {
			return fmt.Sprintf("ERROR: %v", err)
		}
		return "OK"
	case parser.LOAD:
		err := kvStore.LoadFromFile(cmd.Key)
		if err != nil {
			return fmt.Sprintf("ERROR: %v", err)
		}
		return "OK"
	case parser.EXIT:
		return "Goodbye!"
	default:
		return "Unknown command"
	}
}
