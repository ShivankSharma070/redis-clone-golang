package main

import (
	"fmt"
	"log"
	"testing"
)

func TestProtocol(t *testing.T) {
	raw := "set key value\r\n"
	command, err := parseCommand(raw)
	if err != nil {
		log.Fatal("Error ", err)
	}
	fmt.Println(command)
}
