package client

import (
	"testing"
	"context"
	"fmt"
)

func TestNewClient(t *testing.T) {
	t.Log("Starting the test...")
	c, err := New("localhost:5001")
	if err != nil {
		t.Error("Error Creating a client", "err", err)
	}
	for i := range 10 {

		err := c.Set(context.Background(), fmt.Sprintf("name_%d", i), fmt.Sprintf("Shivank_%d", i))
		if err != nil {
			t.Error("Client err in set", "err", err)
		}

		err = c.Get(context.Background(), fmt.Sprintf("name_%d", i))
		if err != nil {
			t.Error("Client err in get", "err", err)
		}

	}
}
