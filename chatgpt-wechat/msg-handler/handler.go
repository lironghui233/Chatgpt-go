package msg_handler

import (
	"chatgpt-wechat/pkg/config"
	"chatgpt-wechat/pkg/log"
	chatgpt_service "chatgpt-wechat/services/chatgpt-service"
	chatgpt_service_proto "chatgpt-wechat/services/chatgpt-service/proto"
	services_client "chatgpt-wechat/services/client"
	"context"
	"strconv"
	"strings"

	"github.com/eatmoreapple/openwechat"
)

type msgHandler struct {
}

func NewMsgHandler() *msgHandler {
	return &msgHandler{}
}

func (mh *msgHandler) TextHandler(ctx *openwechat.MessageContext) {
	var user *openwechat.User
	var group *openwechat.Group
	var content string
	var err error
	var groupID string

	content = ctx.Content
	user, err = ctx.Sender()
	if err != nil {
		return
	}
	//群聊
	if ctx.IsSendByGroup() {
		group = &openwechat.Group{user}
		if !ctx.IsAt() {
			return
		}
		user, err = ctx.SenderInGroup()
		if err != nil {
			log.Error(err)
			return
		}
		content = strings.TrimSpace(strings.ReplaceAll(content, "@"+user.Self().NickName, ""))
		groupID = group.ID()
	}
	replyText, err := generateChatCompletion(user.UserName, groupID, strconv.FormatInt(user.Self().ID(), 10), content)
	if err != nil {
		log.Error(err)
	}
	if ctx.IsSendByGroup() {
		replyText = "@" + user.NickName + " " + replyText
	}
	ctx.ReplyText(replyText)
}
func generateChatCompletion(userID, groupID, endpointAccount, content string) (string, error) {
	cnf := config.GetConf()
	chatGPTServiceClientPool := chatgpt_service.GetChatGPTServiceClientPool()
	conn := chatGPTServiceClientPool.Get()
	defer chatGPTServiceClientPool.Put(conn)

	client := chatgpt_service_proto.NewChatGPTClient(conn)
	ctx1 := context.Background()
	ctx1 = services_client.AppendBearerTokenToContext(ctx1, cnf.DependOnServices.ChatGPTService.AccessToken)
	in := &chatgpt_service_proto.ChatCompletionRequest{
		Id:              userID,
		GroupId:         groupID,
		Message:         content,
		Endpoint:        chatgpt_service_proto.ChatEndpoint_WECHAT,
		EnterpriseId:    cnf.Enterprise.Id,
		EnableContext:   cnf.Chat.EnableContext,
		EndpointAccount: endpointAccount,
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
	res, err := client.ChatCompletion(ctx1, in)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return res.Choices[0].Message.Content, nil
}
