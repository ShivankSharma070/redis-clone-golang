package main

import (
	"bytes"
	"context"
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

func (p *Peer) Write(data string) {
	buf := &bytes.Buffer{}
	rw := resp.NewWriter(buf)
	rw.WriteSimpleString(data)
	_, err := p.conn.Write(buf.Bytes())
	if err != nil {
		slog.Error("Error sending data to client", "err", err, "data", data)
	}
}

func (p *Peer) WriteMap(m map[string]string) error {
	buf := &bytes.Buffer{}
	buf.WriteString("%" + fmt.Sprintf("%d\r\n", len(m)))
	rw := resp.NewWriter(buf)
	for k, v := range m {
		rw.WriteSimpleString(k)
		rw.WriteSimpleString(v)
	}
	_, err := p.conn.Write(buf.Bytes())
	return err
}

func (p *Peer) readLoop(ctx context.Context) error {
	errChan := make(chan error)

	runLoop := func() {
		for {
			command, err := parseCommand(p.conn, p)
			if err != nil {
				errChan <- err
				return
			}
			p.msgCh <- Message{
				data: command,
				peer: p,
			}
		}
	}
	go runLoop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errChan:
			if errors.Is(err, QUIT) || errors.Is(err, net.ErrClosed) {
				return nil
			}
			slog.Error("peer read error", "err", err)
			go runLoop()
		}
	}
}
