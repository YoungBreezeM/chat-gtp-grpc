package client

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
	address = "localhost:8080"
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
				Content: "how to use money",
			},
		},
		ConversationId:  "",
		ParentMessageId: "",
		Token:           "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ik1UaEVOVUpHTkVNMVFURTRNMEZCTWpkQ05UZzVNRFUxUlRVd1FVSkRNRU13UmtGRVFrRXpSZyJ9.eyJodHRwczovL2FwaS5vcGVuYWkuY29tL3Byb2ZpbGUiOnsiZW1haWwiOiI5NDA2OTU4MzZAcXEuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWV9LCJodHRwczovL2FwaS5vcGVuYWkuY29tL2F1dGgiOnsidXNlcl9pZCI6InVzZXItU0dFQkEwT1p3QkFCM2FINUg1NWVDRWF4In0sImlzcyI6Imh0dHBzOi8vYXV0aDAub3BlbmFpLmNvbS8iLCJzdWIiOiJ3aW5kb3dzbGl2ZXwzYjgxYjRiNzYwYzExYWRmIiwiYXVkIjpbImh0dHBzOi8vYXBpLm9wZW5haS5jb20vdjEiLCJodHRwczovL29wZW5haS5vcGVuYWkuYXV0aDBhcHAuY29tL3VzZXJpbmZvIl0sImlhdCI6MTY5MTQ3ODI1MCwiZXhwIjoxNjkyNjg3ODUwLCJhenAiOiJUZEpJY2JlMTZXb1RIdE45NW55eXdoNUU0eU9vNkl0RyIsInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwgbW9kZWwucmVhZCBtb2RlbC5yZXF1ZXN0IG9yZ2FuaXphdGlvbi5yZWFkIG9yZ2FuaXphdGlvbi53cml0ZSBvZmZsaW5lX2FjY2VzcyJ9.A7v5z3cb2uxuAc1jkpxET_EOLnTpZU6nWf4spDclnXLXoHfb6qUxVyyNfY_6RwqUjXvdj6RBn3t46AS8PUDcXE5f0VgRBj1wiKJNcE-5ZoJTQLof7eWGSc65j2a99mnotkar6pfDFFetH3__aeM0I7fAytMWVrx-aWjHJx5H4LYDEoq1uf4WKmgbkjnjJTOwHBWYpu0HCwjAWpSodgEefM7usGsNNzh5qyYO42Uk2lNPEEXP4UHAkPzNJeg-K0xbjeECaVBQuvpB2psYK4z5oXRxTbFel3L365kH_KKOhxjyMXTfXuOq5rQ2nzk8utfZwnaeVhQkqhDYMMVt7PE1NA",
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
