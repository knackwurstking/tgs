all: init build

EXTENSION_TAGS := stats,ip,journal,opmanga,pgpress

clean:
	git clean -xfd

init:
	go mod tidy -v

run:
	go run --tags=${EXTENSION_TAGS} -v ./cmd/tgs

build:
	go build --tags=${EXTENSION_TAGS} -v -o ./bin/tgs ./cmd/tgs
