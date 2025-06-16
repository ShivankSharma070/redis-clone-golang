package client

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	Addr string
	Conn net.Conn
}

func New(remoteAddr string) (*Client, error) {
	myconn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		return nil, err
	}
	return &Client{
		Addr: remoteAddr,
		Conn: myconn,
	}, nil
}

func (c *Client) Set(ctx context.Context, key, value string) error {
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	wr.WriteArray([]resp.Value{resp.StringValue("set"), resp.StringValue(key), resp.StringValue(value)})
	_, err := c.Conn.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write command: %w", err)
	}

	b := make([]byte, 1024)
	_, err = c.Conn.Read(b)
	if err != nil {
		return fmt.Errorf("Error Reading from connection err : %s", err.Error())
	}

	fmt.Println(string(b))
	return err
}
func (c *Client) Get(ctx context.Context, key string) error {

	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	wr.WriteArray([]resp.Value{resp.StringValue("get"), resp.StringValue(key)})
	_, err := c.Conn.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write command: %w", err)

	}

	b := make([]byte, 1024)
	_, err = c.Conn.Read(b)
	if err != nil {
		return fmt.Errorf("Error Reading from connection err : %s", err.Error())
	}
	fmt.Println(string(b))

	return err
}
