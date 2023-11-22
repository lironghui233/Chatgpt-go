package chatgpt

import (
	"chatgpt-qq/config"
	"chatgpt-qq/log"
	chatgpt_service "chatgpt-qq/services/chatgpt-service"
	chatgpt_service_proto "chatgpt-qq/services/chatgpt-service/proto"
	chatgpt_service_client "chatgpt-qq/services/client"
	"context"
	"strconv"
)

func GenerateText(text string, useContext bool, userId, groupId, selfId int64) (string, error) {
	cnf := config.GetConf()
	enterpriseId := cnf.Enterprise.Id
	accessToken := cnf.ChatGPTService.AccessToken
	chatGPTServiceClientPool := chatgpt_service.GetChatGPTServiceClientPool()
	conn := chatGPTServiceClientPool.Get()
	defer chatGPTServiceClientPool.Put(conn)
	ctx := context.Background()
	ctx = chatgpt_service_client.AppendBearerTokenToContext(ctx, accessToken)
	in := &chatgpt_service_proto.ChatCompletionRequest{
		Id:              strconv.FormatInt(userId, 10),
		GroupId:         strconv.FormatInt(groupId, 10),
		Message:         text,
		Endpoint:        chatgpt_service_proto.ChatEndpoint_QQ,
		EnterpriseId:    enterpriseId,
		EnableContext:   useContext,
		EndpointAccount: strconv.FormatInt(selfId, 10),
		ChatParam: &chatgpt_service_proto.ChatParam{
			Model:             cnf.Chat.Model,
			BotDesc:           cnf.Chat.BotDesc,
			ContextLen:        int32(cnf.Chat.ContextLen),
			MinResponseTokens: int32(cnf.Chat.MinResponseTokens),
			ContextTTL:        int32(cnf.Chat.ContextTTL),
			Temperature:       cnf.Chat.Temperature,
			PresencePenalty:   cnf.Chat.PresencePenalty,
			FrequencyPenalty:  cnf.Chat.FrequencyPenalty,
			TopP:              cnf.Chat.TopP,
			MaxTokens:         int32(cnf.Chat.MaxTokens),
		},
	}
	client := chatgpt_service_proto.NewChatGPTClient(conn)
	res, err := client.ChatCompletion(ctx, in)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return res.Choices[0].Message.Content, nil
}
