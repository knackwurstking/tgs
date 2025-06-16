all: init build

clean:
	git clean -xfd

init:
	go mod tidy -v

run:
	go run -v ./cmd/tgs-server

build:
	go build -v -o ./bin/tgs-server ./cmd/tgs-server
