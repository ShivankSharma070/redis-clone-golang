package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/tidwall/resp"
)

const (
	CommandSet = "set"
)

type Command interface {
}

type SetCommand struct {
	key, value string
}

func parseCommand(raw string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(raw))
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch v.Type() {
		case resp.Array:
			switch v.Array()[0].String() {
			case CommandSet:
				if len(v.Array()) != 3 {
					return nil, fmt.Errorf("Not enough number of argument in the command")
				}
				return SetCommand{
					key:   v.Array()[1].String(),
					value: v.Array()[2].String(),
				}, nil
			default:
				return nil, fmt.Errorf("Unkown type of command.")
			}
		default:
			return nil, fmt.Errorf("Unkown type of command.")
		}
	}
	return nil, nil
}
