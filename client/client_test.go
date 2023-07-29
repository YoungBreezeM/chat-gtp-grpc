package main

import (
	"log"
	"sync"
	"testing"

	pb "cgg/api/pb"

	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
)

func Chat() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewChatGTPServiceClient(conn)

	// Contact the server and print out its response.
	chatId := uuid.NewString()
	r, err := c.Chat(context.Background(), &pb.ChatRequest{
		ChatId: chatId,
		Messages: []*pb.ChatRequestMessage{
			{
				Role:    "user",
				Content: "hello",
			},
		},
		Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ik1UaEVOVUpHTkVNMVFURTRNMEZCTWpkQ05UZzVNRFUxUlRVd1FVSkRNRU13UmtGRVFrRXpSZyJ9.eyJodHRwczovL2FwaS5vcGVuYWkuY29tL3Byb2ZpbGUiOnsiZW1haWwiOiI5NDA2OTU4MzZAcXEuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWV9LCJodHRwczovL2FwaS5vcGVuYWkuY29tL2F1dGgiOnsidXNlcl9pZCI6InVzZXItU0dFQkEwT1p3QkFCM2FINUg1NWVDRWF4In0sImlzcyI6Imh0dHBzOi8vYXV0aDAub3BlbmFpLmNvbS8iLCJzdWIiOiJ3aW5kb3dzbGl2ZXwzYjgxYjRiNzYwYzExYWRmIiwiYXVkIjpbImh0dHBzOi8vYXBpLm9wZW5haS5jb20vdjEiLCJodHRwczovL29wZW5haS5vcGVuYWkuYXV0aDBhcHAuY29tL3VzZXJpbmZvIl0sImlhdCI6MTY5MDU5NTU0NiwiZXhwIjoxNjkxODA1MTQ2LCJhenAiOiJUZEpJY2JlMTZXb1RIdE45NW55eXdoNUU0eU9vNkl0RyIsInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwgbW9kZWwucmVhZCBtb2RlbC5yZXF1ZXN0IG9yZ2FuaXphdGlvbi5yZWFkIG9yZ2FuaXphdGlvbi53cml0ZSBvZmZsaW5lX2FjY2VzcyJ9.LyEZqKgiUVaj7i_AuhwoeE_5QXTbEI3Y7hO45iAxp3kDrIeT_45ujpe2X5TIm9DX6sPmanOEsSJd0bkuAiU76sC6C81BkOPBW9e-js53XkkFfvhj5nk_VeolpyGf3gA3izJHHu-9ba42OMPZc3R_qQLl3u5KVaK-mGvxrO6cLdEVGcLtPgnv-iehnZJ1F3JqHCS3Dn9IUyoc42DKb0yS6QDepiz8bd9nsEXp0_Rnhi9xiYQQrKijdtGeT1IMi_hAR4kmFy_EC3O8BcJH2a2bJeB67OCgpL92TU2P3eyEzvVCMLzZLZkHBFPCNUQkYT0elIeJ0llD_-3skfNMVr-Wpg",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	for {
		data, err := r.Recv() //
		if err != nil {
			log.Printf("chat:%s end:%v\n", chatId, err)
			return
		}
		log.Println("resp:", data.Message)
	}
}
func TestChat(t *testing.T) {
	Chat()
}

func TestMultipleClient(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			Chat()
			wg.Done()
		}()
	}
	wg.Wait()
}
