package client

import (
	"bytes"
	"context"
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
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}

	// var buf bytes.Buffer
	// wr := resp.NewWriter(&buf)
	// wr.WriteArray([]resp.Value{resp.StringValue("set"), resp.StringValue("leader"), resp.StringValue("Charlie")})
	// wr.WriteArray([]resp.Value{resp.StringValue("set"), resp.StringValue("follower"), resp.StringValue("Skyler")})
	// fmt.Printf("%s", buf.String())

	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	wr.WriteArray([]resp.Value{resp.StringValue("set"), resp.StringValue(key), resp.StringValue(value)})

	_, err = conn.Write(buf.Bytes())
	return err
}
