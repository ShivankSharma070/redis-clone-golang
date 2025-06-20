package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/tidwall/resp"
)

const (
	CommandSet    = "set"
	CommandGet    = "get"
	CommandHello  = "hello"
	CommandClient = "client"
)

type Command interface {
}

type SetCommand struct {
	key, value []byte
}
type GetCommand struct {
	key []byte
}
type HelloCommand struct {
	value []byte
}
type ClientCommand struct {
	value []byte
}

var QUIT = errors.New("Quit")

// Parse resp command
// Only Set, Get commands are handled
// Hello and client commands are also handled but not implemented, they are just here to bypass official redis client checks.
func parseCommand(reader io.Reader, p *Peer) (Command, error) {
	rd := resp.NewReader(reader)
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			p.delChan <- p
			return nil, fmt.Errorf("End of file %w", QUIT)
		}

		if err != nil {
			return nil, err
		}

		slog.Info("Command received", "cmd", v.String())
		switch v.Type() {
		case resp.Array:
			if len(v.Array()) == 0 {
				return nil, fmt.Errorf("Emply array received")
			}

			switch strings.ToLower(v.Array()[0].String()) {

			// Quit or Exit will end the connection.
			case "quit", "exit":
				return nil, fmt.Errorf("Client requested %w", QUIT)

			case CommandGet:
				if len(v.Array()) != 2 {
					return nil, fmt.Errorf("Bad usage, set accept only 1 argument")
				}
				return GetCommand{
					key: v.Array()[1].Bytes(),
				}, nil

			case CommandSet:
				if len(v.Array()) != 3 {
					return nil, fmt.Errorf("Bad usage, get accept only 1 argument")
				}
				return SetCommand{
					key:   v.Array()[1].Bytes(),
					value: v.Array()[2].Bytes(),
				}, nil

			case CommandHello:
				var value []byte
				if len(v.Array()) > 1 {
					value = v.Array()[1].Bytes()
				}
				return HelloCommand{
					value: []byte(value),
				}, nil

			case CommandClient:
				if len(v.Array()) < 3 {
					return nil, fmt.Errorf("Bad usage, client needs at least 3 arguments.")
				}
				return ClientCommand{
					value: v.Array()[1].Bytes(),
				}, nil

			default:
				return nil, fmt.Errorf("Array, but unkown type .")
			}
		default:
			return nil, fmt.Errorf("Unkown type of command.")
		}
	}
}
