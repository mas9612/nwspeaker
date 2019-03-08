BIN_CRAFTPKT=craftpkt

all: dep test build

build:
	# craftpkt is only for linux because it uses linux specific feature
	GOOS=linux GOARCH=amd64 go build ./cmd/$(BIN_CRAFTPKT)

test:
	go test -v ./...

clean:
	go clean
	rm -f $(BIN_CRAFTPKT)

dep:
	dep ensure
