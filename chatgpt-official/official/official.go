package official

import (
	"bytes"
	"chatgpt-official/pkg/config"
	"chatgpt-official/pkg/log"
	chatgpt_service "chatgpt-official/services/chatgpt-service"
	chatgpt_service_proto "chatgpt-official/services/chatgpt-service/proto"
	services_client "chatgpt-official/services/client"
	crontab "chatgpt-official/services/crontab"
	crontab_client_proto "chatgpt-official/services/crontab/proto"
	"context"
	"crypto/sha1"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func CheckSignature(ctx *gin.Context) {
	cnf := config.GetConf()
	signature := ctx.Query("signature")
	timestamp := ctx.Query("timestamp")
	nonce := ctx.Query("nonce")
	sign := makeSignature(cnf.Official.Token, timestamp, nonce)
	echoStr := ctx.Query("echostr")
	if signature == sign {
		ctx.Data(200, "text/plain;charset=utf-8", []byte(echoStr))
		return
	}
	return
}

type Message struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATA
	FromUserName CDATA
	CreateTime   int64
	MsgType      CDATA
	Content      CDATA
	MsgId        int64
	MsgDataId    int64
	Idx          int
	//图片
	PicUrl  CDATA
	MediaId CDATA
	//语言
	Format      CDATA
	Recognition CDATA
}
type CDATA struct {
	Value string `xml:",cdata"`
}
type AutoReplyTextMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATA
	FromUserName CDATA
	CreateTime   int64
	MsgType      CDATA
	Content      CDATA
}

func ReceiveMessage(ctx *gin.Context) {
	cnf := config.GetConf()
	signature := ctx.Query("signature")
	timestamp := ctx.Query("timestamp")
	nonce := ctx.Query("nonce")
	sign := makeSignature(cnf.Official.Token, timestamp, nonce)
	if signature != sign {
		ctx.Data(200, "text/plain;charset=utf-8", []byte("success"))
		return
	}
	msg := &Message{}
	err := ctx.BindXML(msg)
	if err != nil {
		log.Error(err)
		ctx.Data(200, "text/plain;charset=utf-8", []byte("success"))
		return
	}
	if msg.MsgType.Value != "text" {
		replyMsg := AutoReplyTextMessage{
			ToUserName:   msg.FromUserName,
			FromUserName: msg.ToUserName,
			CreateTime:   time.Now().Unix(),
			MsgType:      CDATA{Value: "text"},
			Content:      CDATA{Value: "抱歉，目前不支持除文本消息以外的其他消息类型"},
		}
		ctx.XML(200, replyMsg)
		return
	}
	ctx.Data(200, "text/plain;charset=utf-8", []byte("success"))
	go func() {
		content, err := generateChatCompletion(msg.FromUserName.Value, msg.ToUserName.Value, msg.Content.Value)
		if err != nil {
			log.Error(err)
			return
		}
		if content == "" {
			return
		}
		sendKfTextMsg(msg.FromUserName.Value, content)
	}()
	return
}

func makeSignature(token, timestamp, nonce string) string {
	sortArr := []string{
		token, timestamp, nonce,
	}
	sort.Strings(sortArr)
	var buffer bytes.Buffer
	for _, value := range sortArr {
		buffer.WriteString(value)
	}
	sha := sha1.New()
	sha.Write(buffer.Bytes()) //使用SHA-1算法对这个字符串进行哈希处理，并将结果格式化为十六进制字符串。
	signature := fmt.Sprintf("%x", sha.Sum(nil))
	return signature
}

type KfTextMessage struct {
	ToUser  string `json:"touser"`
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}
type SendResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func sendKfTextMsg(toUser, content string) {
	accessToken := getAccessToken()
	url := "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + accessToken
	method := "POST"

	replyMsg := &KfTextMessage{
		ToUser:  toUser,
		MsgType: "text",
		Text: struct {
			Content string `json:"content"`
		}{Content: content},
	}
	payloadBytes, err := json.Marshal(replyMsg)
	if err != nil {
		log.Error(err)
		return
	}
	payload := strings.NewReader(string(payloadBytes))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Error(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
		return
	}
	sendRes := &SendResponse{}
	err = json.Unmarshal(body, sendRes)
	if err != nil {
		log.Error(err)
		return
	}
	if sendRes.ErrCode != 0 {
		log.Error(sendRes.ErrMsg)
		return
	}
}

func generateChatCompletion(userID, endpointAccount, content string) (string, error) {
	cnf := config.GetConf()
	chatGPTServiceClientPool := chatgpt_service.GetChatGPTServiceClientPool()
	conn := chatGPTServiceClientPool.Get()
	defer chatGPTServiceClientPool.Put(conn)

	client := chatgpt_service_proto.NewChatGPTClient(conn)
	ctx1 := context.Background()
	ctx1 = services_client.AppendBearerTokenToContext(ctx1, cnf.DependOnServices.ChatGPTService.AccessToken)
	in := &chatgpt_service_proto.ChatCompletionRequest{
		Id:              userID,
		Message:         content,
		Endpoint:        chatgpt_service_proto.ChatEndpoint_WECHATOFFICIAL,
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

func getAccessToken() string {
	cnf := config.GetConf()
	crontabClientPool := crontab.GetCrontabClientPool()
	conn := crontabClientPool.Get()
	defer crontabClientPool.Put(conn)

	client := crontab_client_proto.NewTokenClient(conn)
	in := &crontab_client_proto.TokenRequest{
		Typ: crontab_client_proto.TokenType_WECHATOFFICIAL,
		Id:  cnf.Official.AppId,
		App: "",
	}
	ctx := context.Background()
	ctx = services_client.AppendBearerTokenToContext(ctx, cnf.DependOnServices.Crontab.AccessToken)
	res, err := client.GetToken(ctx, in)
	if err != nil {
		log.Error(err)
		return ""
	}
	return res.AccessToken
}
