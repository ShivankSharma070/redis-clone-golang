package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	listenAddr string
}

type Message struct {
	data Command
	peer *Peer
}

type Server struct {
	Config
	ln          net.Listener
	Peers       map[*Peer]bool // For managing connections
	addPeerChan chan *Peer
	msgChan     chan Message
	delPeerChan chan *Peer
	ctx         context.Context
	kv          *KV // Map to hold key-values pairs
}

func NewServer(ctx context.Context, config Config) *Server {
	return &Server{
		Config:      config,
		Peers:       make(map[*Peer]bool),
		addPeerChan: make(chan *Peer),
		msgChan:     make(chan Message),
		delPeerChan: make(chan *Peer),
		ctx:         ctx,
		kv:          NewKV(),
	}
}

// Listner for tcp connection
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	s.ln = ln

	go s.Loop()
	slog.Info("Server running", "listenAddr", s.listenAddr)
	return s.acceptLoop()
}

// Handle Incoming messages
func (s *Server) handleMessages(msg Message) error {
	switch v := msg.data.(type) {

	case SetCommand:
		err := s.kv.Set(v.key, v.value)
		if err != nil {
			slog.Error("Error setting key value pair", "err", err)
			msg.peer.Write("Error")
		}
		msg.peer.Write("Successfull")

	case GetCommand:
		value, ok := s.kv.Get(v.key)
		if !ok {
			slog.Error("Error in Get command, key not present.")
			msg.peer.Write("Key not present")
		}
		msg.peer.Write(string(value))

	// Below mentioned are command are not implemented by logic, they are just to bypass offical redis client checks.
	// Hello commands expects a map containing server and connection properties.
	case HelloCommand:
		err := msg.peer.WriteMap(map[string]string{"server": "redis"})
		if err != nil {
			msg.peer.Write("Error")
		}

	// Client command expects OK or a err.
	case ClientCommand:
		msg.peer.Write("OK")
	}
	return nil
}

// Manage channels
func (s *Server) Loop() {
	for {
		select {

		case p := <-s.delPeerChan:
			slog.Info("Deleting peer")
			delete(s.Peers, p)

		case msg := <-s.msgChan:
			if err := s.handleMessages(msg); err != nil {
				slog.Error("Raw message error", "err", err)
			}

		case peer := <-s.addPeerChan:
			s.Peers[peer] = true
		}
	}
}

// Accepting Connection
func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept() // Blocking
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				slog.Info("Stop Accepting New requests")
				return nil
			}
			slog.Error("Accept error", "err", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

// Handling Connections
func (s *Server) handleConnection(conn net.Conn) {
	peer := NewPeer(conn, s.msgChan, s.delPeerChan)
	defer func() {
		fmt.Println("Closing the connneciton")
		peer.Write("Closing the connection")
		conn.Close()
	}()
	s.addPeerChan <- peer
	slog.Info("Peer conntected", "connection", conn.LocalAddr())
	peer.readLoop(s.ctx)
}

func main() {
	listenAddr := flag.String("listenAddr", ":5001", "Address to start the server.")
	flag.Parse()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())
	// Connection to server in background
	server := NewServer(ctx, Config{
		listenAddr: *listenAddr,
	})
	time.Sleep(time.Second)
	go func() {
		defer func() { server.ln.Close() }()
		if err := server.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	<-done
	cancel()
	slog.Info("Shutting down server..")
	server.ln.Close()
	time.Sleep(3 * time.Second)
}
