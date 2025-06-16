package client

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
)

// func TestNewClients(t *testing.T) {
// 	t.Log("Creating multiple clients")
// 	wg := sync.WaitGroup{}
// 	for i := range 10 {
//
// 		wg.Add(1)
// 		go func(i int) {
// 			c, err := New("localhost:5001")
// 			if err != nil {
// 				t.Error("Error Creating a client", "err", err)
// 			}
//
// 			defer func() {
// 				c.Conn.Close()
// 				wg.Done()
// 			}()
//
// 			t.Log("Setting value for client", i)
// 			err = c.Set(context.Background(), fmt.Sprintf("name_%d", i), fmt.Sprintf("Shivank_%d", i))
// 			if err != nil {
// 				t.Error("Client err in set", "err", err)
// 			}
//
// 			t.Log("Getting value for client", i)
// 			err = c.Get(context.Background(), fmt.Sprintf("name_%d", i))
// 			if err != nil {
// 				t.Error("Client err in get", "err", err)
// 			}
// 		}(i)
// 	}
//
// 	wg.Wait()
//
// }

func TestRedisOfficialClient(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:5001",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(context.TODO(), "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	// val, err := rdb.Get(context.TODO(), "key").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("key", val)
}
