package main

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net"

	"github.com/tidwall/resp"
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
	buf := &bytes.Buffer{}
	rw := resp.NewWriter(buf)
	rw.WriteBytes(data)
	fmt.Println(buf.String())
	return p.conn.Write(buf.Bytes())
}

func (p *Peer) WriteMap(m map[string]string) {
	buf := &bytes.Buffer{}
	buf.WriteString("%" + fmt.Sprintf("%d\r\n", len(m)))
	rw := resp.NewWriter(buf)
	for k, v := range m {
		rw.WriteSimpleString(k)
		rw.WriteSimpleString(v)
	}
	fmt.Println("This is the map we are sending", buf.String())
	p.conn.Write(buf.Bytes())
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
