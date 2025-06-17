run : build
	@./bin/goredis --listenAddr :5001

build :
	@go build -o bin/goredis

test: 
	@go test -v -count=1

test-client: build
	@{ \
		./bin/goredis --listenAddr :5001 & \
		SERVER_PID=$$!; \
		sleep 1; \
		go test ./client -v -count=1; \
		kill $$SERVER_PID; \
	}
