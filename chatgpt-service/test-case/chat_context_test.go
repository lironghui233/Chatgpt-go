package test_case

import (
	chat_context "chatgpt-service/chatgpt-server/chat-context"
	"chatgpt-service/pkg/config"
	"chatgpt-service/proto"
	"testing"

	"github.com/sashabaranov/go-openai"
)

func TestWebContext(t *testing.T) {
	dataList := []*chat_context.ChatMessage{
		{
			ID: "11111",
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: "你好",
			},
			PID: "",
		},
		{
			ID: "22222",
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "你好，有什么可以帮你吗?",
			},
			PID: "11111",
		},
		{
			ID: "33333",
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: "你好",
			},
			PID: "22222",
		},
		{
			ID: "44444",
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "你好，有什么可以帮你吗?",
			},
			PID: "33333",
		},
		{
			ID: "55555",
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: "你好",
			},
			PID: "44444",
		},
		{
			ID: "66666",
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "你好，有什么可以帮你吗?",
			},
			PID: "55555",
		},
	}
	cache := chat_context.GetCacheContext(proto.ChatEndpoint_WEB)
	for _, item := range dataList {
		err := cache.Set(item.ID, "", proto.ChatEndpoint_WEB, item, config.GetConf().Chat.ContextTTL)
		if err != nil {
			t.Error(err)
			return
		}
	}
	l, err := cache.Get("66666", "", proto.ChatEndpoint_WEB)
	if err != nil {
		t.Error(err)
		return
	}
	cmList, ok := l.([]*chat_context.ChatMessage)
	if !ok {
		t.Error("类型不对")
	}
	if len(cmList) == 0 || len(cmList) > config.GetConf().Chat.ContextLen {
		t.Error("上下文条目不对")
	}
	for _, item := range cmList {
		t.Log(*item)
	}
}
func TestQQContext(t *testing.T) {
	id := "12345678910"
	group := "abcdefghigk"
	dataList := []*chat_context.ChatMessage{
		{
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "你好，有什么可以帮你吗?6",
			},
		},
		{
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: "你好5",
			},
		},
		{
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "你好，有什么可以帮你吗?4",
			},
		},
		{
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: "你好3",
			},
		},
		{
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "你好，有什么可以帮你吗?2",
			},
		},
		{
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: "你好1",
			},
		},
	}
	cache := chat_context.GetCacheContext(proto.ChatEndpoint_QQ)
	err := cache.Set(id, group, proto.ChatEndpoint_QQ, dataList, config.GetConf().Chat.ContextTTL)
	if err != nil {
		t.Error(err)
		return
	}
	l, err := cache.Get(id, group, proto.ChatEndpoint_QQ)
	if err != nil {
		t.Error(err)
		return
	}
	cmList, ok := l.([]*chat_context.ChatMessage)
	if !ok {
		t.Error("类型不对")
	}
	if len(cmList) == 0 || len(cmList) > config.GetConf().Chat.ContextLen {
		t.Error("上下文条目不对")
	}
	for _, item := range cmList {
		t.Log(*item)
	}
}
