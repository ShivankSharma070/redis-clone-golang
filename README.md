# Go Redis Clone

A minimal Redis clone built from scratch in **Go** — just enough to power basic `SET` and `GET` commands, using the same RESP protocol that real Redis speaks. It’s fast, lightweight, and works out of the box with the official [Go Redis client](https://github.com/redis/go-redis). Perfect for learning how Redis works under the hood.

---

##  What It Does

This project mimics a small part of Redis by:

- Accepting TCP connections on port `5001` 
- Parsing commands using the **RESP (REdis Serialization Protocol)**
- Handling the two most basic Redis operations: `SET` and `GET`
- Supporting concurrent connections
- Being fully compatible with Redis clients (tested with Go Redis)
- Running in Docker or natively with a Makefile

---
## ✨ Features

- ✅ `SET` and `GET` support with in-memory storage
-  RESP protocol parsing and encoding (bulk strings, arrays, simple strings, etc.)
- Works with official Redis clients
- Concurrent connections handling.

---

## ️ Setup

### 1. Clone the repo

```bash
git clone https://github.com/ShivankSharma070/redis-clone-golang.git
cd redis-clone-golang
````

---

### 2. Build and run (native)

```bash
make build   # builds the binary
make run     # runs the server (defaults to localhost:5001)
```


To specify the port for redis server set `PORT` while running your make command.

```bash
make run PORT=":3000" # ':' is required.
```


---

### 3. Run with Docker

If you prefer containers:

```bash
docker build . -t go-redis-clone
docker run -p 5001:5001 go-redis-clone 
```

To specify the port for redis server user `--listenAddr` flag :

```bash
docker run -p 3000:3000 go-redis-clone --listenAddr :3000
```

---

##  Using the Go Redis Client

Once the server is running, you can interact with it using the official Go Redis client:

```go
import (
	"fmt"
    "github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:             "localhost:5001",
		DisableIndentity: true, 
	})
	err := rdb.Set(context.Background(), "foo", "bar", 0).Err()
	if err != nil {
		panic(err)
	}
	val, err := rdb.Get(context.TODO(), "foo").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("foo", val)
}
```

---

##  The RESP Protocol (Short Summary)

RESP is the protocol that Redis uses to talk to clients. It's super simple and human-readable:

* `*` means an array
* `$` means a bulk string
* `+` means a simple string
* `-` means an error
* `:` means an integer

Example raw `SET` command over TCP:

```
*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
```

This server parses and handles this, just like Redis would.

---

##  Makefile Commands

| Command            | Description                         |
| ------------------ | ----------------------------------- |
| `make build`       | Compile the server binary           |
| `make run`         | Run the server locally              |
| `make test`        | Test server with multiple Clients   |
| `make test-client` | Test Golang's Official Redis client |
| `make clean`       | Remove binaries and build artifacts |

---

## ⚠️ Limitations

This is a basic prototype meant for learning. Some things it **does not** do yet:

* No persistence (all data is in-memory)
* No support for other Redis commands
* No clustering, pub/sub, or authentication
* Limited error handling

---

##  Why I Built This

Redis is a brilliant piece of software — simple in concept, blazing fast in practice. By building a tiny clone from scratch, I wanted to deeply understand:

* TCP servers in Go
* Protocol design and parsing
* Memory stores and command dispatching
* How real clients interact with Redis over the wire

This project helped me gain that insight, and I hope it helps you too.

---

##  Acknowledgments
This project was originally inspired by the following video tutorial:

[Building a Redis Clone in Go](https://www.youtube.com/watch?v=LMrxfWB6sbQ&t=2551s)

All core logic and protocol implementation were based on that walkthrough. I've since modified, extended, and containerized the project to suit my own learning and understanding of how Redis works internally.

Big thanks to the original creator for making such a helpful resource!

* [Redis Protocol Spec](https://redis.io/docs/latest/develop/reference/protocol-spec/)
* [Go Redis Client](https://github.com/redis/go-redis)
