all: server client

server:
	go build server.go shared.go

client:
	go build client.go shared.go
