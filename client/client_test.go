package client

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/tidwall/resp"
)

func TestNewClients(t *testing.T) {
	t.Log("Creating multiple clients")
	wg := sync.WaitGroup{}
	for i := range 10 {

		wg.Add(1)
		go func(i int) {
			c, err := New("localhost:5001")
			if err != nil {
				t.Error("Error Creating a client", "err", err)
			}

			defer func() {
				c.Conn.Close()
				wg.Done()
			}()

			t.Log("Setting value for client", i)
			err = c.Set(context.Background(), fmt.Sprintf("name_%d", i), fmt.Sprintf("Shivank_%d", i))
			if err != nil {
				t.Error("Client err in set", "err", err)
			}

			t.Log("Getting value for client", i)
			err = c.Get(context.Background(), fmt.Sprintf("name_%d", i))
			if err != nil {
				t.Error("Client err in get", "err", err)
			}
		}(i)
	}

	wg.Wait()

}

func TestRedisOfficialClient(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:             "localhost:5001",
		Password:         "", // no password set
		DB:               0,  // use default DB
		DisableIndentity: true,
	})

	fmt.Println("Sending set command")
	err := rdb.Set(context.Background(), "name", "shivank", 0).Err()
	if err != nil {
		fmt.Println("This is the error", err.Error())
		panic(err)
	}

	fmt.Println("Sending get command")
	val, err := rdb.Get(context.TODO(), "name").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("name", val)
}

func TestMapWriting(t *testing.T) {
	m := map[string]string{"foo": "bar"}
	buf := &bytes.Buffer{}
	buf.WriteString("%" + fmt.Sprintf("%d\r\n", len(m)))
	rw := resp.NewWriter(buf)
	for k, v := range m {
		rw.WriteSimpleString(k)
		rw.WriteSimpleString(v)
	}
	fmt.Println("This is the map we are sending", buf.String())
}
