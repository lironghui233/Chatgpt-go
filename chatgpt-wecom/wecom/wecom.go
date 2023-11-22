package wecom

import (
	"chatgpt-wecom/pkg/config"
	"chatgpt-wecom/pkg/log"
	chatgpt_service "chatgpt-wecom/services/chatgpt-service"
	chatgpt_service_proto "chatgpt-wecom/services/chatgpt-service/proto"
	services_client "chatgpt-wecom/services/client"
	crontab "chatgpt-wecom/services/crontab"
	crontab_client_proto "chatgpt-wecom/services/crontab/proto"
	"chatgpt-wecom/wxbizmsgcrypt"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckSignature(ctx *gin.Context) {
	cnf := config.GetConf()
	verifyMsgSign := ctx.Query("msg_signature")
	verifyTimestamp := ctx.Query("timestamp")
	verifyNonce := ctx.Query("nonce")
	verifyEchoStr := ctx.Query("echostr")

	wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(cnf.Wecom.Token, cnf.Wecom.EncodingAesKey, cnf.Wecom.CorpId, wxbizmsgcrypt.XmlType)

	echoStr, cryptErr := wxcpt.VerifyURL(verifyMsgSign, verifyTimestamp, verifyNonce, verifyEchoStr)
	if cryptErr != nil {
		return
	}
	ctx.Data(200, "text/plain;charset=utf-8", echoStr)
	return
}

type CDATA struct {
	Value string `xml:",cdata"`
}
type Message struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName CDATA
	CreateTime int64
	MsgType    CDATA
	Event      CDATA
	Token      CDATA
	OpenKfId   CDATA
}

func ReceiveMessage(ctx *gin.Context) {
	cnf := config.GetConf()
	reqMsgSign := ctx.Query("msg_signature")
	reqTimestamp := ctx.Query("timestamp")
	reqNonce := ctx.Query("nonce")
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Error(err)
		return
	}

	wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(cnf.Wecom.Token, cnf.Wecom.EncodingAesKey, cnf.Wecom.CorpId, wxbizmsgcrypt.XmlType)
	msgBytes, cryptErr := wxcpt.DecryptMsg(reqMsgSign, reqTimestamp, reqNonce, body)
	if cryptErr != nil {
		log.Error(cryptErr)
		return
	}
	msg := &Message{}
	err = xml.Unmarshal(msgBytes, msg)
	if err != nil {
		log.Error(err)
		return
	}
	//调用客服信息获取接口，获取客户发来的客服消息
	cursor := ""
getNext:
	syncMsgRes, err := getMsg(cursor, msg.Token.Value, msg.OpenKfId.Value)
	if err != nil {
		log.Error(err)
		return
	}
	if syncMsgRes.HasMore == 1 {
		cursor = syncMsgRes.NextCursor
		goto getNext
	}
	//处理消息
	go handleMsg(syncMsgRes)
	return
}

type SyncMsgRequest struct {
	Cursor      string `json:"cursor"`
	Token       string `json:"token"`
	Limit       int    `json:"limit"`
	VoiceFormat int    `json:"voice_format"`
	OpenKfId    string `json:"open_kfid"`
}
type SyncMsgResponse struct {
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	NextCursor string `json:"next_cursor"`
	HasMore    int    `json:"has_more"`
	MsgList    []Msg  `json:"msg_list"`
}
type Msg struct {
	MsgId          string `json:"msgid"`
	OpenKfId       string `json:"open_kfid"`
	ExternalUserId string `json:"external_userid"`
	SendTime       int64  `json:"send_time"`
	Origin         int    `json:"origin"`
	ServicerUserId string `json:"servicer_user_id"`
	MsgType        string `json:"msgtype"`
	Text           Text   `json:"text"`
	Image          Image  `json:"image"`
	Event          Event  `json:"event"`
}
type Text struct {
	Content string `json:"content"`
	MenuId  string `json:"menu_id"`
}
type Image struct {
	MediaId string `json:"media_id"`
}
type Event struct {
	EventType      string `json:"event_type"`
	OpenKfId       string `json:"open_kfid"`
	ExternalUserId string `json:"external_userid"`
	Scene          string `json:"scene"`
	SceneParam     string `json:"scene_param"`
	WelcomeCode    string `json:"welcome_code"`
}

func getMsg(cursor, token, openKfId string) (*SyncMsgResponse, error) {
	accessToken := getAccessToken()
	url := "https://qyapi.weixin.qq.com/cgi-bin/kf/sync_msg?access_token=" + accessToken
	method := "POST"

	reqPayload := &SyncMsgRequest{
		Cursor:   cursor,
		Token:    token,
		OpenKfId: openKfId,
	}
	payloadBytes, err := json.Marshal(reqPayload)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	payload := strings.NewReader(string(payloadBytes))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	syncMsgResponse := &SyncMsgResponse{}
	err = json.Unmarshal(body, syncMsgResponse)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return syncMsgResponse, err
}

type SendMsgRequest struct {
	ToUser   string `json:"touser"`
	OpenKfId string `json:"open_kfid"`
	MsgId    string `json:"msgid"`
	MsgType  string `json:"msgtype"`
	Text     Text   `json:"text"`
	Image    Image  `json:"image"`
}
type SendMsgResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	MsgId   string `json:"msgid"`
}

func sendKfTextMsgToUser(toUser, openKfId, msgId, content string) error {
	accessToken := getAccessToken()
	url := "https://qyapi.weixin.qq.com/cgi-bin/kf/send_msg?access_token=" + accessToken
	method := "POST"

	reqPayload := &SendMsgRequest{
		ToUser:   toUser,
		OpenKfId: openKfId,
		MsgType:  "text",
		MsgId:    msgId,
		Text:     Text{Content: content},
	}
	payloadBytes, err := json.Marshal(reqPayload)
	if err != nil {
		log.Error(err)
		return err
	}
	payload := strings.NewReader(string(payloadBytes))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.Error(err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
		return err
	}
	sendMsgResponse := &SendMsgResponse{}
	err = json.Unmarshal(body, sendMsgResponse)
	if err != nil {
		log.Error(err)
		return err
	}
	if sendMsgResponse.ErrCode != 0 {
		log.Error(sendMsgResponse.ErrMsg)
		return errors.New(sendMsgResponse.ErrMsg)
	}
	return nil
}
func handleMsg(response *SyncMsgResponse) {
	if response.ErrCode != 0 {
		log.Error(response.ErrMsg)
		return
	}
	if len(response.MsgList) < 1 {
		return
	}
	current := response.MsgList[len(response.MsgList)-1]
	if current.MsgType != "text" {
		return
	}
	content, err := generateChatCompletion(current.ExternalUserId, current.OpenKfId, current.Text.Content)
	if err != nil {
		return
	}
	//发送客服消息
	err = sendKfTextMsgToUser(current.ExternalUserId, current.OpenKfId, current.MsgId, content)
	if err != nil {
		log.Error(err)
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
		Typ: crontab_client_proto.TokenType_WECOM,
		Id:  cnf.Wecom.CorpId,
		App: "kf",
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
