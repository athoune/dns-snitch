build:
	go build .
	upx dns-snitch

test:
	go test ./...
