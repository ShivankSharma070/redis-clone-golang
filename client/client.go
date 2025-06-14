package client

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	addr string
}

func New(remoteAddr string) *Client {

	return &Client{
		addr: remoteAddr,
	}
}

func (c *Client) Set(ctx context.Context, key, value string) error {
	// FIX: Optmize this, avoid making connection for every query
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	wr.WriteArray([]resp.Value{resp.StringValue("set"), resp.StringValue(key), resp.StringValue(value)})
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write command: %w", err)
	}

	return err
}
func (c *Client) Get(ctx context.Context, key string) error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	wr.WriteArray([]resp.Value{resp.StringValue("get"), resp.StringValue(key)})
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write command: %w", err)

	}

	b := make([]byte, 1024)
	_, err = conn.Read(b)
	if err != nil {
		return fmt.Errorf("Error Reading from connection err : %s", err.Error())
	}
	fmt.Println(string(b))

	return err
}
