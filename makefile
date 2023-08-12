.PHONY:build

all:build 

build:
	GOOS=linux GOARCH=amd64 go build -o build/server -v -ldflags="-s -w" cmd/main.go && upx -9 build/server

build-client:
	GOOS=linux GOARCH=amd64 go build -o build/client -v -ldflags="-s -w" client/main.go && upx -9 build/client

clear:
	rm -r build/*