all: server client

clean:
	rm server client

fmt: 
	go fmt *.go
run:
	go run server.go client.go shared.go

server:
	go build server.go shared.go client.go

chat:
	go build client.go shared.go
