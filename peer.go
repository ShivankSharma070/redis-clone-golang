package main

import (
	"net"
	"errors"
	"log/slog"
)

type Peer struct {
	conn  net.Conn
	msgCh chan Message
}

func NewPeer(conn net.Conn, msgChan chan Message) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgChan,
	}
}

func (p *Peer) Write(data []byte) (int, error) {
	return p.conn.Write(append(data, byte('\n')))
}

func (p *Peer) readLoop() error {
	for {
		command, err := parseCommand(p.conn)

		if errors.Is(err, QUIT) {
			return err
		}

		if err != nil {
			slog.Error("Peer read error", "err", err)
			continue
		}
		p.msgCh <- Message{
			data: command,
			peer: p,
		}
	}
}
