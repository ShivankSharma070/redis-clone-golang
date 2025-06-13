build :
	@go build -o bin/goredis

run : build
	@./bin/goredis

