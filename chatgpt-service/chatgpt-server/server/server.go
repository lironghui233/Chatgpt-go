package server

import (
	chat_context "chatgpt-service/chatgpt-server/chat-context"
	"chatgpt-service/pkg/config"
	"chatgpt-service/pkg/log"
	"chatgpt-service/proto"
	chatgpt_data "chatgpt-service/services/chatgpt-data"
	chatgpt_data_proto "chatgpt-service/services/chatgpt-data/proto"
	services_client "chatgpt-service/services/client"
	"chatgpt-service/services/keywords"
	keywords_proto "chatgpt-service/services/keywords/proto"
	sensitive_words "chatgpt-service/services/sensitive-words"
	sensitive_words_proto "chatgpt-service/services/sensitive-words/proto"
	"chatgpt-service/services/tokenizer"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

const ChatPrimedTokens = 2

type chatGPTServer struct {
	proto.UnimplementedChatGPTServer
	config *config.Config
	log    log.ILogger
}

func NewChatGPTServer(config *config.Config, log log.ILogger) proto.ChatGPTServer {
	return &chatGPTServer{
		config: config,
		log:    log,
	}
}

type chatGPTAPP struct {
	config config.Config
	log    log.ILogger
}

func (s *chatGPTServer) getChatGPTAPP(in *proto.ChatCompletionRequest) *chatGPTAPP {
	//注意只是取指针的值进行赋值，因此不会修改原server的conf值
	conf := *s.config
	if in.ChatParam != nil {
		if in.ChatParam.Model != "" {
			conf.Chat.Model = in.ChatParam.Model
		}
		conf.Chat.TopP = in.ChatParam.TopP
		conf.Chat.FrequencyPenalty = in.ChatParam.FrequencyPenalty
		conf.Chat.PresencePenalty = in.ChatParam.PresencePenalty
		conf.Chat.Temperature = in.ChatParam.Temperature
		conf.Chat.BotDesc = in.ChatParam.BotDesc

		if conf.Chat.MaxTokens != 0 {
			conf.Chat.MaxTokens = int(in.ChatParam.MaxTokens)
		}
		if conf.Chat.ContextTTL != 0 {
			conf.Chat.ContextTTL = int(in.ChatParam.ContextTTL)
		}
		if conf.Chat.MinResponseTokens != 0 {
			conf.Chat.MinResponseTokens = int(in.ChatParam.MinResponseTokens)
		}
		if conf.Chat.ContextLen != 0 {
			conf.Chat.ContextLen = int(in.ChatParam.ContextLen)
		}
	}
	app := &chatGPTAPP{
		log:    s.log,
		config: conf, //这里conf的内容已经被替换为传过来的参数，跟server的conf内容不一样，server的内容并没有被修改
	}
	return app
}

func (s *chatGPTServer) ChatCompletion(ctx context.Context, in *proto.ChatCompletionRequest) (*proto.ChatCompletionResponse, error) {
	app := s.getChatGPTAPP(in)
	//敏感词过滤
	ok, msg, err := app.sensitiveWords(in)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	if !ok {
		res := app.buildChatCompletionResponse(msg)
		return res, nil
	}
	//提取关键词查找
	keywords := app.keywords(in)

	client := app.getChatGPTClient()
	contextList, tokensNum, currTokensNum, currMessage, req, err := app.buildChatCompletionRequest(in, false)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	res := &proto.ChatCompletionResponse{}
	bytes, err := json.Marshal(resp)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	err = jsonpb.UnmarshalString(string(bytes), res)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	//保存上下文
	go func() {
		reqContext := &chat_context.ChatMessage{
			Message:   currMessage,
			TokensNum: currTokensNum,
		}
		resContext := &chat_context.ChatMessage{
			ID:        uuid.New().String(),
			Message:   resp.Choices[0].Message,
			TokensNum: resp.Usage.CompletionTokens,
		}
		if in.Endpoint == proto.ChatEndpoint_WEB {
			reqContext.ID = in.Id
			reqContext.PID = in.Pid

			if resp.ID != "" {
				resContext.ID = resp.ID
			}
			resContext.PID = reqContext.ID
		}
		err := app.saveContext(in, reqContext, resContext, contextList)
		if err != nil {
			s.log.Error(err)
			return
		}
	}()
	//调用数据服务
	go func() {
		err := app.saveData(in, keywords, currTokensNum, resp.Choices[0].Message.Content, resp.Usage.CompletionTokens, tokensNum)
		if err != nil {
			s.log.Error(err)
		}
	}()
	return res, err
}
func (s *chatGPTServer) ChatCompletionStream(in *proto.ChatCompletionRequest, stream proto.ChatGPT_ChatCompletionStreamServer) error {
	app := s.getChatGPTAPP(in)
	//敏感词过滤
	ok, msg, err := app.sensitiveWords(in)
	if err != nil {
		s.log.Error(err)
		return err
	}
	if !ok {
		resId := uuid.New().String()
		startRes := app.buildChatCompletionStreamResponse(resId, "", "")
		endRes := app.buildChatCompletionStreamResponse(resId, "", "stop")
		err = stream.Send(startRes)
		if err != nil {
			s.log.Error(err)
			return err
		}
		resList := app.buildChatCompletionStreamResponseList(resId, msg)
		for _, res := range resList {
			err = stream.Send(res)
			if err != nil {
				s.log.Error(err)
				return err
			}
		}
		err = stream.Send(endRes)
		if err != nil {
			s.log.Error(err)
			return err
		}
		return nil
	}
	//提取关键词查找
	keywords := app.keywords(in)

	client := app.getChatGPTClient()
	contextList, tokensNum, currTokensNum, currMessage, req, err := app.buildChatCompletionRequest(in, true)
	if err != nil {
		s.log.Error(err)
		return err
	}
	chatStream, err := client.CreateChatCompletionStream(stream.Context(), req)
	if err != nil {
		s.log.Error(err)
		return err
	}
	defer chatStream.Close()
	completionContent := ""
	resultID := ""
	for {
		resp, err := chatStream.Recv()
		if err != nil && err != io.EOF {
			s.log.Error(err)
			return err
		}
		if err == io.EOF {
			break
		}
		if resultID == "" {
			resultID = resp.ID
		}
		completionContent += resp.Choices[0].Delta.Content
		res := &proto.ChatCompletionStreamResponse{}

		bytes, err := json.Marshal(resp)
		if err != nil {
			s.log.Error(err)
			return err
		}
		err = jsonpb.UnmarshalString(string(bytes), res)
		if err != nil {
			s.log.Error(err)
			return err
		}
		err = stream.Send(res)
		if err != nil {
			s.log.Error(err)
			return err
		}
	}

	//保存上下文
	resultCompletion := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: completionContent,
	}
	resultTokensNum, err := tokenizer.GetTokenCount(resultCompletion, s.config.Chat.Model)
	if err != nil {
		s.log.Error(err)
		return err
	}

	go func() {
		reqContext := &chat_context.ChatMessage{
			Message:   currMessage,
			TokensNum: currTokensNum,
		}
		resContext := &chat_context.ChatMessage{
			ID: uuid.New().String(),
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: completionContent,
			},
			TokensNum: resultTokensNum,
		}
		if in.Endpoint == proto.ChatEndpoint_WEB {
			reqContext.ID = in.Id
			reqContext.PID = in.Pid

			if resultID != "" {
				resContext.ID = resultID
			}
			resContext.PID = reqContext.ID
		}
		err := app.saveContext(in, reqContext, resContext, contextList)
		if err != nil {
			s.log.Error(err)
			return
		}
	}()

	// 调用数据服务
	go func() {
		err := app.saveData(in, keywords, currTokensNum, completionContent, resultTokensNum, tokensNum)
		if err != nil {
			s.log.Error(err)
		}
	}()
	return nil
}

func (app *chatGPTAPP) getChatGPTClient() *openai.Client {
	conf := app.config
	accessToken := conf.Chat.APIKey
	config := openai.DefaultConfig(accessToken)
	config.BaseURL = conf.Chat.BaseURL
	client := openai.NewClientWithConfig(config)
	return client
}
func (app *chatGPTAPP) buildChatCompletionRequest(in *proto.ChatCompletionRequest, stream bool) (contextList []*chat_context.ChatMessage, tokensNum, currTokensNum int, currMessage openai.ChatCompletionMessage, req openai.ChatCompletionRequest, err error) {
	conf := app.config
	currMessage = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: in.Message,
	}
	req = openai.ChatCompletionRequest{
		Model: conf.Chat.Model,
		Messages: []openai.ChatCompletionMessage{
			currMessage,
		},
		// 表示当前请求的回复，最大token数
		MaxTokens:        conf.Chat.MinResponseTokens,
		Temperature:      conf.Chat.Temperature,
		TopP:             conf.Chat.TopP,
		PresencePenalty:  conf.Chat.PresencePenalty,
		FrequencyPenalty: conf.Chat.FrequencyPenalty,
		Stream:           stream,
	}
	contextList = make([]*chat_context.ChatMessage, 0)
	var value interface{}
	if in.EnableContext {
		cache := chat_context.GetCacheContext(in.Endpoint)
		cacheID := in.Id
		if in.Endpoint == proto.ChatEndpoint_WEB {
			cacheID = in.Pid
		}

		value, err = cache.Get(cacheID, in.GroupId, in.Endpoint)
		if err != nil {
			app.log.Error(err)
			return
		}
		contextList = value.([]*chat_context.ChatMessage)
	}
	tokensNum, currTokensNum, req.Messages, err = app.buildMessages(contextList, currMessage)
	if err != nil {
		app.log.Error(err)
		return
	}
	req.MaxTokens = conf.Chat.MaxTokens - tokensNum
	return
}

func (app *chatGPTAPP) buildMessages(contextList []*chat_context.ChatMessage, currMessage openai.ChatCompletionMessage) (tokensNum int, currTokensNum int, messages []openai.ChatCompletionMessage, err error) {
	conf := app.config
	var sysMessage openai.ChatCompletionMessage
	messages = []openai.ChatCompletionMessage{currMessage}
	currTokensNum, err = tokenizer.GetTokenCount(currMessage, conf.Chat.Model)
	if err != nil {
		app.log.Error(err)
		return
	}
	botTokens := 0
	if conf.Chat.BotDesc != "" {
		sysMessage = openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: conf.Chat.BotDesc,
		}
		botTokens, err = tokenizer.GetTokenCount(sysMessage, conf.Chat.Model)
		if err != nil {
			log.Error(err)
			return
		}
	}
	if currTokensNum > conf.Chat.MaxTokens-conf.Chat.MinResponseTokens-botTokens {
		err = errors.New(fmt.Sprintf("上下文合计tokens最大%d，回复保留tokens数%d，ai特征使用tokens %d，剩余可用tokens %d，当前消息tokens %d", conf.Chat.MaxTokens, conf.Chat.MinResponseTokens, botTokens, conf.Chat.MaxTokens-conf.Chat.MinResponseTokens-botTokens, currTokensNum))
		log.Error(err)
		return
	}
	tokensNum = currTokensNum + botTokens + ChatPrimedTokens
	for _, item := range contextList {
		if tokensNum+item.TokensNum > conf.Chat.MaxTokens-conf.Chat.MinResponseTokens {
			break
		}
		messages = append(messages, item.Message)
		tokensNum += item.TokensNum + ChatPrimedTokens
	}
	//反转messages列表
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	if botTokens > 0 {
		messages = append([]openai.ChatCompletionMessage{sysMessage}, messages...)
	}
	return
}

func (app *chatGPTAPP) saveContext(in *proto.ChatCompletionRequest, reqContext, resContext *chat_context.ChatMessage, contextList []*chat_context.ChatMessage) (err error) {
	cache := chat_context.GetCacheContext(in.Endpoint)
	if in.Endpoint == proto.ChatEndpoint_WEB {
		err = cache.Set(reqContext.ID, in.GroupId, in.Endpoint, reqContext, app.config.Chat.ContextTTL)
		if err != nil {
			app.log.Error(err)
			return
		}
		err = cache.Set(resContext.ID, in.GroupId, in.Endpoint, resContext, app.config.Chat.ContextTTL)
		if err != nil {
			app.log.Error(err)
			return
		}
		return nil
	}
	contextList = append([]*chat_context.ChatMessage{resContext, reqContext}, contextList...)
	if len(contextList) > app.config.Chat.ContextLen {
		contextList = contextList[:app.config.Chat.ContextLen]
	}
	err = cache.Set(in.Id, in.GroupId, in.Endpoint, contextList, app.config.Chat.ContextTTL)
	if err != nil {
		app.log.Error(err)
		return
	}
	return nil
}

func (app *chatGPTAPP) sensitiveWords(in *proto.ChatCompletionRequest) (ok bool, msg string, err error) {
	sensitiveClientPool := sensitive_words.GetSensitiveWordsClientPool()
	conn := sensitiveClientPool.Get()
	defer sensitiveClientPool.Put(conn)

	client := sensitive_words_proto.NewSensitiveWordsClient(conn)
	ctx := context.Background()
	ctx = services_client.AppendBearerTokenToContext(ctx, app.config.DependOnServices.SensitiveWords.AccessToken)
	req := &sensitive_words_proto.ValidateReq{
		Text: in.Message,
	}
	res, err := client.Validate(ctx, req)
	if err != nil {
		log.Error(err)
		return false, "", err
	}
	ok = res.Ok
	if !ok {
		msg = "触及到了知识盲区哦，换个问题再问吧"
	}
	return
}
func (app *chatGPTAPP) keywords(in *proto.ChatCompletionRequest) []string {
	keywordsClientPool := keywords.GetKeywordsClientPool()
	conn := keywordsClientPool.Get()
	defer keywordsClientPool.Put(conn)

	client := keywords_proto.NewKeywordsClient(conn)
	ctx := context.Background()
	ctx = services_client.AppendBearerTokenToContext(ctx, app.config.DependOnServices.Keywords.AccessToken)
	req := &keywords_proto.FindAllReq{
		Text: in.Message,
	}
	res, err := client.FindAll(ctx, req)
	if err != nil {
		app.log.Error(err)
		return []string{}
	}
	return res.Keywords
}
func (app *chatGPTAPP) saveData(in *proto.ChatCompletionRequest, keywords []string, userMsgTokens int, aiMsg string, aiMsgTokens int, reqTokens int) error {
	dataClientPool := chatgpt_data.GetChatGPTDataClientPool()
	conn := dataClientPool.Get()
	defer dataClientPool.Put(conn)

	client := chatgpt_data_proto.NewChatGPTDataClient(conn)
	ctx := context.Background()
	ctx = services_client.AppendBearerTokenToContext(ctx, app.config.DependOnServices.ChatGPTData.AccessToken)
	req := &chatgpt_data_proto.Record{
		UserMsg:         in.Message,
		UserMsgTokens:   int32(userMsgTokens),
		AiMsg:           aiMsg,
		AiMsgTokens:     int32(aiMsgTokens),
		UserMsgKeywords: keywords,
		ReqTokens:       int32(reqTokens),
		CreateAt:        time.Now().Unix(),
		Endpoint:        int32(in.Endpoint),
		EnterpriseId:    in.EnterpriseId,
		EndpointAccount: in.EndpointAccount,
	}
	if in.Endpoint != proto.ChatEndpoint_WEB {
		req.Account = in.Id
		req.GroupId = in.GroupId
	}
	_, err := client.AddRecord(ctx, req)
	return err
}

func (app *chatGPTAPP) buildChatCompletionResponse(msg string) *proto.ChatCompletionResponse {
	res := &proto.ChatCompletionResponse{
		Id:      uuid.New().String(),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   app.config.Chat.Model,
		Choices: []*proto.ChatCompletionChoice{
			{
				Message: &proto.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleAssistant,
					Content: msg,
				},
				FinishReason: "stop",
			},
		},
		Usage: &proto.Usage{
			PromptTokens:     0,
			CompletionTokens: 0,
			TotalTokens:      0,
		},
	}
	return res
}

func (app *chatGPTAPP) buildChatCompletionStreamResponseList(id, msg string) []*proto.ChatCompletionStreamResponse {
	list := make([]*proto.ChatCompletionStreamResponse, 0)
	for _, delta := range msg {
		list = append(list, app.buildChatCompletionStreamResponse(id, string(delta), ""))
	}
	return list
}
func (app *chatGPTAPP) buildChatCompletionStreamResponse(id, delta, finishReason string) *proto.ChatCompletionStreamResponse {
	res := &proto.ChatCompletionStreamResponse{
		Id:      id,
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   app.config.Chat.Model,
		Choices: []*proto.ChatCompletionStreamChoice{
			{
				Index: 0,
				Delta: &proto.ChatCompletionStreamChoiceDelta{
					Content: delta,
					Role:    openai.ChatMessageRoleAssistant,
				},
				FinishReason: finishReason,
			},
		},
	}
	return res
}
