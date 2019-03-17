BIN_CRAFTPKT=craftpkt

.PHONY: all clean

all: test dep build

build:
	# craftpkt is only for linux because it uses linux specific feature
	GOOS=linux GOARCH=amd64 go build ./cmd/$(BIN_CRAFTPKT)

test:
	docker build -t nwspeaker-test -f Dockerfile.test .
	docker run --rm nwspeaker-test

clean:
	go clean
	rm -f $(BIN_CRAFTPKT)

dep:
	dep ensure
