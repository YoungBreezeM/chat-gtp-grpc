
build:
	GOOS=linux GOARCH=amd64 go build -o server -v -ldflags="-s -w" cmd/main.go && upx -9 server