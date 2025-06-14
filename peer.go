package main

import (
	"net"
)

type Peer struct {
	conn  net.Conn
	msgCh chan []byte
}

func NewPeer(conn net.Conn, msgChan chan []byte) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgChan,
	}
}

func (p *Peer) readLoop() error {
	// Reading data form a connection
	buf := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			return err
		}
		msgBuf := make([]byte, n)
		copy(msgBuf, buf)
		p.msgCh <- msgBuf
	}
}
