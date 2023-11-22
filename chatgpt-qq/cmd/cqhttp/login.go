package cqhttp

import (
	"chatgpt-qq/log"
	"encoding/json"
	"os"
	"regexp"
	"strconv"
	"time"

	"chatgpt-qq/config"

	"github.com/gorilla/websocket"
)

type Bot struct {
	// 机器人QQ号
	QQ   int64
	Conn *websocket.Conn
}

func NewBot() *Bot {
	return &Bot{}
}

func Run() {
	var bot = NewBot()
	logger := log.NewLogger()
	logger.SetOutput(os.Stderr)
	logger.SetLevel("info")
	for i := 1; ; i++ {
		logger.InfoF("第%d次尝试连接%s中...\n", i, config.Cfg.CqHttp.WebSocket)
		var err error
		bot.Conn, _, err = websocket.DefaultDialer.Dial(config.Cfg.CqHttp.WebSocket, nil)
		if err != nil {
			logger.InfoF("连接失败, 5秒后重试:%v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		logger.Info("CqHttp连接成功!")
		go bot.Read()
		break
	}
}

func (bot *Bot) Read() {
	conn := bot.Conn
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Error(err)
			break
		}
		if msgType == websocket.CloseMessage {
			conn.Close()
			log.Info("连接关闭")
			break
		}
		var rcvMsg RcvMsg
		err = json.Unmarshal(msg, &rcvMsg)
		if err != nil {
			log.Error(err)
			continue
		}
		if bot.QQ == 0 && rcvMsg.SelfId != 0 {
			bot.QQ = rcvMsg.SelfId
		}
		//处理收到的消息
		if rcvMsg.PostType == "message" {
			// 消息预处理Parser
			isAt, err := regexp.MatchString(`CQ:at,qq=`+strconv.FormatInt(rcvMsg.SelfId, 10), rcvMsg.RawMessage)
			if err != nil {
				log.Error(err)
				continue
			}
			// 去除消息CQ码
			rcvMsg.Message = regexp.MustCompile(`\[CQ:.*?]`).ReplaceAllString(rcvMsg.Message, "")
			if rcvMsg.Message != " " && rcvMsg.Message != "" && rcvMsg.Message != "   " {
				go bot.HandleMsg(isAt, rcvMsg)
			}
		}
	}
}
