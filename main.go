package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"

	"github.com/ShivankSharma070/redis-clone-go/client"
)

const default_listen_Addr = ":5001"

type Config struct {
	listenAddr string
}
type Server struct {
	Config
	ln          net.Listener
	Peers       map[*Peer]bool // For managing connections
	addPeerChan chan *Peer
	quitChan    chan struct{}
	msgChan     chan []byte

	kv *KV
}

func NewServer(config Config) *Server {
	if len(config.listenAddr) == 0 {
		config.listenAddr = default_listen_Addr
	}

	return &Server{
		Config:      config,
		Peers:       make(map[*Peer]bool),
		addPeerChan: make(chan *Peer),
		quitChan:    make(chan struct{}),
		msgChan:     make(chan []byte),
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
func (s *Server) rawMessageHandler(rawMsg []byte) error {
	command, err := parseCommand(string(rawMsg))
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
		fmt.Println(string(value))
	}

	return nil
}

//Manage connections
func (s *Server) Loop() {
	for {
		select {
		case rawMsg := <-s.msgChan:
			if err := s.rawMessageHandler(rawMsg); err != nil {
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
	server := NewServer(Config{})
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
	}

	client.Get(context.Background(),"name_1")
	select {} // Blocking }
}
