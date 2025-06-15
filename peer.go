package main

import (
	"net"
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
	return p.conn.Write(data)
}
func (p *Peer) readLoop() error {
	for {
		command, err := parseCommand(p.conn)
		if err != nil {
			return err
		}
		p.msgCh <- Message{
			data: command,
			peer: p,
		}
	}
}
