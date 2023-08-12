package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"cgg/internal/di"

	"google.golang.org/grpc"
)

var (
	port = ":8080"
)

func main() {

	serverPort := os.Getenv("SERVER_PORT")
	if len(serverPort) > 0 {
		port = fmt.Sprintf(":%s", serverPort)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	di.Wire(s)

	log.Printf("server start on %s\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
