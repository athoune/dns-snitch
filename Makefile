all: build

build:
	go build .

test:
	go test ./...

clean:
	rm -f dns-snitch
