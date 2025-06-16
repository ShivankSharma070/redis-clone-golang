package main

import (
	"errors"
	"log/slog"
	"net"
)

type Peer struct {
	conn    net.Conn
	msgCh   chan Message
	delChan chan *Peer
}

func NewPeer(conn net.Conn, msgChan chan Message, delChan chan *Peer) *Peer {
	return &Peer{
		conn:    conn,
		msgCh:   msgChan,
		delChan: delChan,
	}
}

func (p *Peer) Write(data []byte) (int, error) {
	return p.conn.Write(append(data, byte('\n')))
}

func (p *Peer) readLoop() error {
	for {
		command, err := parseCommand(p.conn, p)

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
