package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/zahidhasann88/kvstore/parser"
	"github.com/zahidhasann88/kvstore/server"
	"github.com/zahidhasann88/kvstore/store"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "server":
			fmt.Println("Starting KV Store TCP Server on :8080...")
			server.StartTCPServer(":8080")
			return
		case "client":
			fmt.Println("Connecting to KV Store server...")
			server.StartTCPClient("localhost:8080")
			return
		}
	}

	kvStore := store.NewStore()
	defer kvStore.Close()

	fmt.Println("Simple Key-Value Store CLI")
	fmt.Println("Commands: SET key value, GET key, DEL key, SAVE filename, LOAD filename, EXIT")
	fmt.Println("Start TCP server with: go run . server")
	fmt.Println("Connect to server with: go run . client")
	fmt.Println("Note: Use quotes for values with spaces: SET key \"hello world\"")

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

		cmd, err := parser.ParseCommand(input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		result := executeCommand(kvStore, cmd)
		fmt.Println(result)

		if cmd.Type == parser.EXIT {
			break
		}
	}
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
		return fmt.Sprintf("\"%s\"", value)
	case parser.DEL:
		deleted := kvStore.Delete(cmd.Key)
		if deleted {
			return "1"
		}
		return "0"
	case parser.SAVE:
		err := kvStore.SaveToFile(cmd.Key)
		if err != nil {
			return fmt.Sprintf("Error: %v", err)
		}
		return "OK"
	case parser.LOAD:
		err := kvStore.LoadFromFile(cmd.Key)
		if err != nil {
			return fmt.Sprintf("Error: %v", err)
		}
		return "OK"
	case parser.EXIT:
		return "Goodbye!"
	default:
		return "Unknown command"
	}
}
