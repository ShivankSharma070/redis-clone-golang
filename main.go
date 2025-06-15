package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
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

	kv *KV
}

func NewServer(config Config) *Server {
	return &Server{
		Config:      config,
		Peers:       make(map[*Peer]bool),
		addPeerChan: make(chan *Peer),
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
	switch v := msg.data.(type) {
	case SetCommand:
		err := s.kv.Set(v.key, v.value)
		if err != nil {
			slog.Error("Error setting key value pair", "err", err)
		}
		msg.peer.Write([]byte("successfull"))

	case GetCommand:
		value, present := s.kv.Get(v.key)
		if !present {
			return fmt.Errorf("Get Command error, no such key exists")
		}

		msg.peer.Write(value)
	}

	return nil
}

// Manage connections
func (s *Server) Loop() {
	for {
		select {
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
			slog.Error("Accept error", "err", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

// Handling Connections
func (s *Server) handleConnection(conn net.Conn) {
	peer := NewPeer(conn, s.msgChan)
	defer func() { 
		peer.Write([]byte("Closing the connection"))
		conn.Close() 
	}()
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
	time.Sleep(time.Second)
	go func() {
		defer func (){server.ln.Close()}()
		log.Fatal(server.Start())
	}()

	select {} // Blocking
}
