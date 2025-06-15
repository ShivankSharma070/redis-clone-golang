run : build
	@./bin/goredis --listenAddr :5001

build :
	@go build -o bin/goredis

test: 
	@go test -v -count=1
