package service

import (
	"bufio"
	"bytes"
	"cgg/api/pb"
	"cgg/internal/constant/chatgtp"
	"cgg/pkg/tlsc"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	http "github.com/bogdanfinn/fhttp"

	"github.com/google/uuid"
)

type ChatGTP struct {
}

type ChatGPTRequest struct {
	Action                     string           `json:"action"`
	Messages                   []ChatGTPMessage `json:"messages"`
	ParentMessageID            string           `json:"parent_message_id,omitempty"`
	ConversationID             string           `json:"conversation_id,omitempty"`
	Model                      string           `json:"model"`
	HistoryAndTrainingDisabled bool             `json:"history_and_training_disabled"`
}

type ChatGTPAuthor struct {
	Role string `json:"role"`
}

type ChatGTPContent struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type ChatGTPMessage struct {
	ID      uuid.UUID      `json:"id"`
	Author  ChatGTPAuthor  `json:"author"`
	Content ChatGTPContent `json:"content"`
}

func NewChatGPTRequest() ChatGPTRequest {
	return ChatGPTRequest{
		Action:                     "next",
		ParentMessageID:            uuid.NewString(),
		Model:                      "text-davinci-002-render-sha",
		HistoryAndTrainingDisabled: false,
	}
}

func (c *ChatGPTRequest) Send(req *pb.ChatRequest, ch chan string) {

	for _, m := range req.Messages {
		c.Messages = append(c.Messages, ChatGTPMessage{
			ID:      uuid.New(),
			Author:  ChatGTPAuthor{Role: m.Role},
			Content: ChatGTPContent{ContentType: "text", Parts: []string{m.Content}},
		})
	}

	client, err := tlsc.NewTLSClient()
	if err != nil {
		log.Println(err.Error())
		return
	}

	//set proxy
	proxy := os.Getenv("PROXY_URL")
	if len(proxy) > 0 {
		if err = client.SetProxy(proxy); err != nil {
			log.Println(err)
			return
		}
	}

	body, err := json.Marshal(c)
	if err != nil {
		log.Println(err)
		return
	}
	request, err := http.NewRequest(http.MethodPost, chatgtp.APIURL, bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
	}
	// Clear cookies
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36")
	request.Header.Set("Accept", "*/*")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", req.Token))
	request.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")

	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()
	//
	if response.StatusCode == 200 {
		reader := bufio.NewReader(response.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					ch <- io.EOF.Error()
					break
				}
				log.Println("Failed to read response:", err)
				return
			}
			//
			if len(line) > 1 {
				ch <- line
			}
		}

	} else {
		log.Println(response.Status)
	}

}

var chatting = false

func (s *ChatGTP) Chat(req *pb.ChatRequest, res pb.ChatGTPService_ChatServer) (err error) {
	log.Printf("Received: %s", req.ChatId)

	if chatting {
		if err = res.Send(&pb.ChatReply{Message: chatgtp.CHATTING}); err != nil {
			log.Println(err)
		}
		return
	}

	c := make(chan string)
	cg := NewChatGPTRequest()

	chatting = true
	go cg.Send(req, c)

	for {
		msg := <-c
		if msg == "EOF" {
			chatting = false
			break
		}
		if err = res.Send(&pb.ChatReply{Message: msg}); err != nil {
			log.Println(err)
		}
	}
	return nil
}
