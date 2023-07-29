build:
	GOOS=linux GOARCH=amd64 go build -o build/server -v -ldflags="-s -w" cmd/main.go && upx -9 build/server