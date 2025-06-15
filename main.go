package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"
	"flag"

	"github.com/ShivankSharma070/redis-clone-go/client"
)


type Config struct {
	listenAddr string
}

type Message struct {
	data []byte
	peer *Peer
}

type Server struct {
	Config
	ln          net.Listener
	Peers       map[*Peer]bool // For managing connections
	addPeerChan chan *Peer
	quitChan    chan struct{}
	msgChan     chan Message

	kv *KV
}

func NewServer(config Config) *Server {
	return &Server{
		Config:      config,
		Peers:       make(map[*Peer]bool),
		addPeerChan: make(chan *Peer),
		quitChan:    make(chan struct{}),
		msgChan:     make(chan Message),
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
	command, err := parseCommand(string(msg.data))
	if err != nil {
		return err
	}

	switch v := command.(type) {
	case SetCommand:
		s.kv.Set(v.key, v.value)
	case GetCommand:
		value, present := s.kv.Get(v.key)
		if !present { 
			return fmt.Errorf("Get Command error, no such key exists")
		}

		msg.peer.Write(append(value, byte('\n')))
	}

	return nil
}

//Manage connections
func (s *Server) Loop() {
	for {
		select {
		case msg := <-s.msgChan:
			if err := s.handleMessages(msg); err != nil {
				slog.Error("Raw message error", "err", err)
			}
		case <-s.quitChan:
			return
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
			slog.Error("Accept error", "err", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

// Handling Connections
func (s *Server) handleConnection(conn net.Conn) {
	peer := NewPeer(conn, s.msgChan)
	s.addPeerChan <- peer
	slog.Info("Peer conntected", "connection", conn.LocalAddr())
	err := peer.readLoop()
	if err != nil {
		slog.Error("Peer read error", "err", err)
	}
}

func main() {
	listenAddr := flag.String("listenAddr", ":5001", "Address to start the server.")
	flag.Parse()

	// Connection to server in background
	server := NewServer(Config{
		listenAddr: *listenAddr,
	})
	go func() {
		log.Fatal(server.Start())
	}()
	time.Sleep(2 * time.Second)

	client := client.New("localhost:5001")
	for i := range 10 {
		
		err := client.Set(context.Background(), fmt.Sprintf("name_%d", i), fmt.Sprintf("Shivank_%d", i))
		if err != nil {
			slog.Error("Client err in set", "err", err)
		}
		
		err= client.Get(context.Background(),fmt.Sprintf("name_%d",i))
		if err != nil {
			slog.Error("Client err in get", "err", err)
		}

	}

	select {} // Blocking }
}
