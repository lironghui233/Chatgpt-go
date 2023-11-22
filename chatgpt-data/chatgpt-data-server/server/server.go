package server

import (
	"chatgpt-data/chatgpt-data-server/data"
	"chatgpt-data/pkg/config"
	"chatgpt-data/pkg/log"
	"chatgpt-data/proto"
	"context"
	"time"
)

type ChatGPTDataServer struct {
	//内嵌结构，拥有这个内嵌结构的所有方法。通过内嵌去实现AddRecord()方法
	proto.UnimplementedChatGPTDataServer
	//配置信息
	config *config.Config
	//log
	log log.ILogger
	//数据访问层
	chatRecordsData data.IChatRecordsData
}

func NewChatGPTDataServer(conf *config.Config, log log.ILogger, chatRecordsData data.IChatRecordsData) proto.ChatGPTDataServer {
	return &ChatGPTDataServer{
		config:          conf,
		log:             log,
		chatRecordsData: chatRecordsData,
	}
}

func (s *ChatGPTDataServer) AddRecord(ctx context.Context, in *proto.Record) (*proto.RecordRes, error) {
	cr := &data.ChatRecord{}
	cr.Account = in.Account
	cr.UserMsgKeywords = in.UserMsgKeywords
	cr.GroupID = in.GroupId
	cr.AIMsg = in.AiMsg
	cr.UserMsg = in.UserMsg
	cr.AIMsgTokens = int(in.AiMsgTokens)
	cr.UserMsgTokens = int(in.UserMsgTokens)
	cr.CreateAt = in.CreateAt
	cr.ReqTokens = int(in.ReqTokens)
	cr.EndpointAccount = in.EndpointAccount
	cr.Endpoint = int(in.Endpoint)
	cr.EnterpriseId = in.EnterpriseId
	if cr.CreateAt == 0 {
		cr.CreateAt = time.Now().Unix()
	}
	err := s.chatRecordsData.Add(cr)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	out := &proto.RecordRes{
		Id: cr.ID,
	}
	return out, err
}
