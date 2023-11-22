package test_case

import (
	"context"
	"io"
	"testing"

	"github.com/sashabaranov/go-openai"
)

// testing.B用于定义一个性能测试
// testing.T用于定义一个常规的测试
func BenchmarkProxyChatCompletion(b *testing.B) {
	accessToken := conf.GetString("http.access_token")
	config := openai.DefaultConfig(accessToken)
	config.BaseURL = *addr
	client := openai.NewClientWithConfig(config)
	req := openai.ChatCompletionRequest{
		Model: conf.GetString("chat.model"),
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "你好",
			},
		},
		MaxTokens:        512,
		Temperature:      float32(conf.GetFloat64("chat.temperature")),
		TopP:             float32(conf.GetFloat64("chat.top_p")),
		PresencePenalty:  float32(conf.GetFloat64("chat.presence_penalty")),
		FrequencyPenalty: float32(conf.GetFloat64("chat.frequency_penalty")),
	}
	req = setBotDesc(req, "你是一个AI助手，我需要你模拟一名资深的软件工程师来回答我的问题")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.CreateChatCompletion(context.Background(), req)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkProxyChatCompletionStream(b *testing.B) {
	accessToken := conf.GetString("http.access_token")
	config := openai.DefaultConfig(accessToken)
	config.BaseURL = *addr
	client := openai.NewClientWithConfig(config)
	req := openai.ChatCompletionRequest{
		Model: conf.GetString("chat.model"),
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "你好",
			},
		},
		MaxTokens:        512,
		Temperature:      float32(conf.GetFloat64("chat.temperature")),
		TopP:             float32(conf.GetFloat64("chat.top_p")),
		PresencePenalty:  float32(conf.GetFloat64("chat.presence_penalty")),
		FrequencyPenalty: float32(conf.GetFloat64("chat.frequency_penalty")),
		Stream:           true,
	}
	req = setBotDesc(req, "你是一个AI助手，我需要你模拟一名资深的软件工程师来回答我的问题")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stream, err := client.CreateChatCompletionStream(context.Background(), req)
		if err != nil {
			b.Error(err)
			return
		}
		var contentList = make([]string, 0)
		for {
			resp, err := stream.Recv()
			if err != nil && err != io.EOF {
				b.Error(err)
				return
			}
			if err == io.EOF {
				break
			}
			contentList = append(contentList, resp.Choices[0].Delta.Content)
		}
		stream.Close()
	}

}

func BenchmarkProxyChatCompletionParallel(b *testing.B) {
	accessToken := conf.GetString("http.access_token")
	config := openai.DefaultConfig(accessToken)
	config.BaseURL = *addr
	client := openai.NewClientWithConfig(config)
	req := openai.ChatCompletionRequest{
		Model: conf.GetString("chat.model"),
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "你好",
			},
		},
		MaxTokens:        512,
		Temperature:      float32(conf.GetFloat64("chat.temperature")),
		TopP:             float32(conf.GetFloat64("chat.top_p")),
		PresencePenalty:  float32(conf.GetFloat64("chat.presence_penalty")),
		FrequencyPenalty: float32(conf.GetFloat64("chat.frequency_penalty")),
	}
	req = setBotDesc(req, "你是一个AI助手，我需要你模拟一名资深的软件工程师来回答我的问题")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.CreateChatCompletion(context.Background(), req)
			if err != nil {
				b.Error(err)
				return
			}
		}
	})
}
