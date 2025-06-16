package main

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/tidwall/resp"
)

const (
	CommandSet = "set"
	CommandGet = "get"
)

type Command interface {
}

type SetCommand struct {
	key, value []byte
}
type GetCommand struct {
	key []byte
}

var QUIT = errors.New("Quit")

func parseCommand(reader io.Reader, p *Peer) (Command, error) {
	rd := resp.NewReader(reader)
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			p.delChan <- p
			break
		}

		if err != nil {
			return nil, err
		}

		// TODO: Handle all types
		fmt.Println("Command got", v.String())
		switch v.Type() {
		case resp.Array:
			if len(v.Array()) == 0 {
				return nil, fmt.Errorf("Emply array received")
			}

			switch strings.ToLower(v.Array()[0].String()) {
			case "quit", "exit":
				return nil, fmt.Errorf("Client requested %w", QUIT)
			case CommandGet:
				if len(v.Array()) != 2 {
					return nil, fmt.Errorf("Not enough number of argument in the get command")
				}
				return GetCommand{
					key: v.Array()[1].Bytes(),
				}, nil
			case CommandSet:
				if len(v.Array()) != 3 {
					return nil, fmt.Errorf("Not enough number of argument in the set command")
				}
				return SetCommand{
					key:   v.Array()[1].Bytes(),
					value: v.Array()[2].Bytes(),
				}, nil
			default:
				return nil, fmt.Errorf("Array, but unkown type of array.")
			}
		default:
			return nil, fmt.Errorf("Unkown type of command.")
		}
	}
	return nil, nil
}
