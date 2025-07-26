package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CommandType int

const (
	SET CommandType = iota
	GET
	DEL
	SAVE
	LOAD
	EXIT
)

type Command struct {
	Type  CommandType
	Key   string
	Value string
	TTL   time.Duration
}

func ParseCommand(input string) (*Command, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	cmd := strings.ToUpper(parts[0])

	switch cmd {
	case "SET":
		if len(parts) < 3 {
			return nil, fmt.Errorf("SET requires key and value")
		}

		valueStart := 2
		valueEnd := len(parts)

		if len(parts) >= 5 && strings.ToUpper(parts[len(parts)-2]) == "EX" {
			seconds, err := strconv.Atoi(parts[len(parts)-1])
			if err != nil {
				return nil, fmt.Errorf("invalid TTL value")
			}
			valueEnd = len(parts) - 2

			command := &Command{
				Type:  SET,
				Key:   parts[1],
				Value: cleanQuotes(strings.Join(parts[valueStart:valueEnd], " ")),
				TTL:   time.Duration(seconds) * time.Second,
			}
			return command, nil
		}

		command := &Command{
			Type:  SET,
			Key:   parts[1],
			Value: cleanQuotes(strings.Join(parts[valueStart:valueEnd], " ")),
		}

		return command, nil

	case "GET":
		if len(parts) != 2 {
			return nil, fmt.Errorf("GET requires exactly one key")
		}
		return &Command{
			Type: GET,
			Key:  parts[1],
		}, nil

	case "DEL":
		if len(parts) != 2 {
			return nil, fmt.Errorf("DEL requires exactly one key")
		}
		return &Command{
			Type: DEL,
			Key:  parts[1],
		}, nil

	case "SAVE":
		if len(parts) != 2 {
			return nil, fmt.Errorf("SAVE requires filename")
		}
		return &Command{
			Type: SAVE,
			Key:  parts[1],
		}, nil

	case "LOAD":
		if len(parts) != 2 {
			return nil, fmt.Errorf("LOAD requires filename")
		}
		return &Command{
			Type: LOAD,
			Key:  parts[1],
		}, nil

	case "EXIT", "QUIT":
		return &Command{
			Type: EXIT,
		}, nil

	default:
		return nil, fmt.Errorf("unknown command: %s", cmd)
	}
}

func cleanQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
