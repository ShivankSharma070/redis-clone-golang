package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/ShivankSharma070/redis-clone-go/client"
)

// Test server with multipe clients
func TestNewServerWithClients(t *testing.T) {

	listenAddr := flag.String("listenAddr", ":5001", "Address to start the server.")
	flag.Parse()

	// Connection to server in background
	server := NewServer(context.Background(), Config{
		listenAddr: *listenAddr,
	})
	time.Sleep(time.Second)

	go func() {
		defer func() { server.ln.Close() }()
		log.Fatal(server.Start())
	}()

	wg := sync.WaitGroup{}
	for i := range 10 {

		wg.Add(1)
		go func(i int) {
			c, err := client.New("localhost:5001")
			if err != nil {
				t.Error("Error Creating a client", "err", err)
			}

			defer func() {
				c.Conn.Close()
				wg.Done()
			}()

			t.Log("Setting value for client", i)
			err = c.Set(context.Background(), fmt.Sprintf("foo_%d", i), fmt.Sprintf("bar_%d", i))
			if err != nil {
				t.Error("Client err in set", "err", err)
			}

			t.Log("Getting value for client", i)
			err = c.Get(context.Background(), fmt.Sprintf("foo_%d", i))
			if err != nil {
				t.Error("Client err in get", "err", err)
			}
		}(i)
	}

	wg.Wait()

	// Check if there is any peers left
	slog.Info("Peers still connected", "No of Connected peers", len(server.Peers))

	time.Sleep(time.Second)
	if len(server.Peers) > 0 {
		panic("Peers left.")
	}

}
