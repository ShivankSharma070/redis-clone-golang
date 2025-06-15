package client

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func TestNewClients(t *testing.T) {
	t.Log("Creating multiple clients")
	wg := sync.WaitGroup{}
	for i := range 10 {
		
		wg.Add(1)
		go func(i int) {
			c, err := New("localhost:5001")
			if err != nil {
				t.Error("Error Creating a client", "err", err)
			}

			defer func(){
				c.conn.Close()
				wg.Done()
			}()

			t.Log("Setting value for client", i)
			err = c.Set(context.Background(), fmt.Sprintf("name_%d", i), fmt.Sprintf("Shivank_%d", i))
			if err != nil {
				t.Error("Client err in set", "err", err)
			}

			t.Log("Getting value for client", i)
			err = c.Get(context.Background(), fmt.Sprintf("name_%d", i))
			if err != nil {
				t.Error("Client err in get", "err", err)
			}
		}(i)
	}

	wg.Wait()

}
