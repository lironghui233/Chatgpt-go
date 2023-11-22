package controllers

import (
	"chatgpt-web/pkg/config"
	"chatgpt-web/pkg/log"
	chatgpt_service "chatgpt-web/services/chatgpt-service"
	chatgpt_service_proto "chatgpt-web/services/chatgpt-service/proto"
	services_client "chatgpt-web/services/client"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"k8s.io/klog/v2"
)

const (
	ChatPrimedTokens = 2
)

type ChatService struct {
	config *config.Config
	log    log.ILogger
}

type ChatMessageRequest struct {
	Prompt  string                    `json:"prompt"`
	Options ChatMessageRequestOptions `json:"options"`
}

type ChatMessageRequestOptions struct {
	Name            string `json:"name"`
	ParentMessageId string `json:"parentMessageId"`
}

type ChatMessage struct {
	ID              string                                              `json:"id"`
	Text            string                                              `json:"text"`
	Role            string                                              `json:"role"`
	Name            string                                              `json:"name"`
	Delta           string                                              `json:"delta"`
	Detail          *chatgpt_service_proto.ChatCompletionStreamResponse `json:"detail"`
	TokenCount      int                                                 `json:"tokenCount"`
	ParentMessageId string                                              `json:"parentMessageId"`
}

const ChatMessageRoleAssistant = "assistant"

func NewChatService(config *config.Config, log log.ILogger) (*ChatService, error) {
	chat := ChatService{
		log:    log,
		config: config,
	}
	return &chat, nil
}

func (chat *ChatService) ChatProcess(ctx *gin.Context) {
	payload := ChatMessageRequest{}
	if err := ctx.BindJSON(&payload); err != nil {
		chat.log.Error(err)
		ctx.JSON(200, gin.H{
			"status":  "Fail",
			"message": fmt.Sprintf("%v", err),
			"data":    nil,
		})
		return
	}

	messageID := uuid.New().String()
	result := ChatMessage{
		ID:              uuid.New().String(),
		Role:            ChatMessageRoleAssistant,
		Text:            "",
		ParentMessageId: messageID,
	}
	chatGPTServiceClientPool := chatgpt_service.GetChatGPTServiceClientPool()
	conn := chatGPTServiceClientPool.Get()
	defer chatGPTServiceClientPool.Put(conn)

	client := chatgpt_service_proto.NewChatGPTClient(conn)
	ctx1 := context.Background()
	ctx1 = services_client.AppendBearerTokenToContext(ctx1, chat.config.DependOnServices.ChatGPTService.AccessToken)
	in := &chatgpt_service_proto.ChatCompletionRequest{
		Id:              messageID,
		Message:         payload.Prompt,
		Pid:             payload.Options.ParentMessageId,
		Endpoint:        chatgpt_service_proto.ChatEndpoint_WEB,
		EnterpriseId:    chat.config.Enterprise.Id,
		EnableContext:   false,
		EndpointAccount: chat.config.Enterprise.Id,
		ChatParam: &chatgpt_service_proto.ChatParam{
			Model:             chat.config.Chat.Model,
			BotDesc:           chat.config.Chat.BotDesc,
			ContextLen:        int32(chat.config.Chat.ContextLen),
			MinResponseTokens: int32(chat.config.Chat.MinResponseTokens),
			ContextTTL:        int32(chat.config.Chat.ContextTTL),
			Temperature:       chat.config.Chat.Temperature,
			PresencePenalty:   chat.config.Chat.PresencePenalty,
			FrequencyPenalty:  chat.config.Chat.FrequencyPenalty,
			TopP:              chat.config.Chat.TopP,
			MaxTokens:         int32(chat.config.Chat.MaxTokens),
		},
	}
	if payload.Options.ParentMessageId != "" {
		in.EnableContext = true
	}
	stream, err := client.ChatCompletionStream(ctx1, in)
	if err != nil {
		chat.log.Error(err)
		ctx.JSON(200, gin.H{
			"status":  "Fail",
			"message": fmt.Sprintf("%v", err),
			"data":    nil,
		})
		return
	}
	defer stream.CloseSend()

	firstChunk := true
	ctx.Header("Content-type", "application/octet-stream")
	for {
		rsp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return
		}
		if err != nil {
			chat.log.Error(err)
			ctx.JSON(200, gin.H{"status": "Fail", "message": fmt.Sprintf("OpenAI Event Error %v", err), "data": nil})
			return
		}
		if rsp.Id != "" {
			result.ID = rsp.Id
		}
		if len(rsp.Choices) > 0 {
			content := rsp.Choices[0].Delta.Content
			result.Delta = content
			if len(content) > 0 {
				result.Text += content
			}
			result.Detail = rsp
		}

		bts, err := json.Marshal(result)
		if err != nil {
			chat.log.Error(err)
			ctx.JSON(200, gin.H{"status": "Fail", "message": fmt.Sprintf("OpenAI Event Marshal Error %v", err), "data": nil})
			return
		}

		if !firstChunk {
			ctx.Writer.Write([]byte("\n"))
		} else {
			firstChunk = false
		}

		if _, err := ctx.Writer.Write(bts); err != nil {
			klog.Error(err)
			return
		}
		ctx.Writer.Flush()
	}
}
