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
	// Reading data form a connection
	buf := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			return err
		}
		msgBuf := make([]byte, n)
		copy(msgBuf, buf)
		p.msgCh <- Message{
			data: msgBuf,
			peer: p,
		}
	}
}
