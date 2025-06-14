all: init build

init:
	go mod tidy -v

run:
	go run -v ./cmd/tgs-server

build:
	go build -v -o ./bin/tgs-server ./cmd/tgs-server
