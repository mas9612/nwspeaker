BIN_NWSPEAKER=nwspeaker
BIN_ARPSPEAKER=arpspeaker

.PHONY: all clean

all: dep test build

build:
	GOOS=linux GOARCH=amd64 go build ./cmd/$(BIN_NWSPEAKER)
	GOOS=linux GOARCH=amd64 go build ./cmd/$(BIN_ARPSPEAKER)

test:
	docker build -t nwspeaker-test -f Dockerfile.test .
	docker run --rm nwspeaker-test
	# remove dangling image
	docker image prune -f

clean:
	go clean
	rm -f $(BIN_ARPSPEAKER)
	docker image prune -f

dep:
	dep ensure
