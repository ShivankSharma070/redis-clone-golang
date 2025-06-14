package main

import (
	"bytes"
	"fmt"
	"io"

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
			if len(v.Array()) == 0 {
				return nil, fmt.Errorf("Emply array received")
			}

			switch v.Array()[0].String() {
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
				return nil, fmt.Errorf("Unkown type of command.")
			}
		default:
			return nil, fmt.Errorf("Unkown type of command.")
		}
	}
	return nil, nil
}
