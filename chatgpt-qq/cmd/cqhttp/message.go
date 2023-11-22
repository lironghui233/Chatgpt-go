package cqhttp

import (
	"chatgpt-qq/cmd/chatgpt"
	"chatgpt-qq/config"
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

type RcvMsg struct {
	PostType    string `json:"post_type"`
	MessageType string `json:"message_type"`
	Time        int64  `json:"time"`
	SelfId      int64  `json:"self_id"`
	SubType     string `json:"sub_type"`
	UserId      int64  `json:"user_id"`
	TargetId    int64  `json:"target_id"`
	Message     string `json:"message"`
	RawMessage  string `json:"raw_message"`
	Font        int    `json:"font"`
	Sender      struct {
		Age      int    `json:"age"`
		Nickname string `json:"nickname"`
		Sex      string `json:"sex"`
		UserId   int64  `json:"user_id"`
	}
	GroupId   int64 `json:"group_id"`
	MessageId int64 `json:"message_id"`
}
type SendMsg struct {
	Action string `json:"action"`
	Params struct {
		UserId  int64  `json:"user_id"`
		GroupId int64  `json:"group_id"`
		Message string `json:"message"`
	} `json:"params"`
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags) //初始化日志格式
}

// HandleMsg 对CqHttp发送的json进行处理
func (bot *Bot) HandleMsg(isAt bool, rcvMsg RcvMsg) {
	// panic处理 暂时无法判断是否生效
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic:", err)
			return
		}
	}()
	switch rcvMsg.MessageType {
	case "private":
		// 包含关键词才触发，运算符优先级 && > ||
		if config.Cfg.CqHttp.UseKeyword && !strings.Contains(rcvMsg.Message, config.Cfg.CqHttp.Keyword) && (config.Cfg.CqHttp.KeywordType == "all" || config.Cfg.CqHttp.KeywordType == "private") || rcvMsg.Sender.UserId == bot.QQ {
			return
		}
		rcvMsg.Message = strings.ReplaceAll(rcvMsg.Message, config.Cfg.CqHttp.Keyword, "")
		msg, err := chatgpt.GenerateText(rcvMsg.Message, config.Cfg.Context.PrivateContext, rcvMsg.UserId, rcvMsg.GroupId, rcvMsg.SelfId)
		if msg != "" {
			err = bot.SendPrivateMsg(rcvMsg.Sender.UserId, "[CQ:reply,id="+strconv.FormatInt(rcvMsg.MessageId, 10)+"]"+msg)
		} else {
			err = bot.SendPrivateMsg(rcvMsg.Sender.UserId, "[CQ:reply,id="+strconv.FormatInt(rcvMsg.MessageId, 10)+"]"+"生成错误！错误信息:\n"+err.Error())
		}
		if err != nil {
			log.Println(err)
		}
	case "group":
		// 群消息@机器人才处理
		if !isAt && config.Cfg.CqHttp.AtOnly || rcvMsg.Sender.UserId == bot.QQ {
			return
		}
		// 检查是否有关键词
		if config.Cfg.CqHttp.UseKeyword && !strings.Contains(rcvMsg.Message, config.Cfg.CqHttp.Keyword) && (config.Cfg.CqHttp.KeywordType == "all" || config.Cfg.CqHttp.KeywordType == "group") {
			return
		}
		rcvMsg.Message = strings.ReplaceAll(rcvMsg.Message, config.Cfg.CqHttp.Keyword, "")

		msg, err := chatgpt.GenerateText(rcvMsg.Message, config.Cfg.Context.GroupContext, rcvMsg.UserId, rcvMsg.GroupId, rcvMsg.SelfId)
		if msg != "" {
			err = bot.SendGroupMsg(rcvMsg.GroupId, "[CQ:reply,id="+strconv.FormatInt(rcvMsg.MessageId, 10)+"]"+msg)
		} else {
			err = bot.SendGroupMsg(rcvMsg.GroupId, "[CQ:reply,id="+strconv.FormatInt(rcvMsg.MessageId, 10)+"]"+"生成错误！错误信息：\n"+err.Error())
		}
		if err != nil {
			log.Println(err)
		}
	}

}

func (bot *Bot) SendPrivateMsg(userId int64, text string) error {
	sendMsg := SendMsg{
		Action: "send_private_msg",
		Params: struct {
			UserId  int64  `json:"user_id"`
			GroupId int64  `json:"group_id"`
			Message string `json:"message"`
		}{UserId: userId, Message: text},
	}

	rawMsg, _ := json.Marshal(sendMsg)
	err := bot.Conn.WriteMessage(1, rawMsg)
	if err != nil {
		return err
	}
	return nil
}
func (bot *Bot) SendGroupMsg(groupId int64, text string) error {
	sendMsg := SendMsg{
		Action: "send_group_msg",
		Params: struct {
			UserId  int64  `json:"user_id"`
			GroupId int64  `json:"group_id"`
			Message string `json:"message"`
		}{GroupId: groupId, Message: text},
	}

	rawMsg, _ := json.Marshal(sendMsg)
	err := bot.Conn.WriteMessage(1, rawMsg)
	if err != nil {
		return err
	}
	return nil
}
