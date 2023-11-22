package test_case

import (
	"chatgpt-data/chatgpt-data-server/data"
	"chatgpt-data/chatgpt-data-server/server"
	"chatgpt-data/pkg/config"
	"chatgpt-data/pkg/db/mysql"
	"chatgpt-data/pkg/log"
	"chatgpt-data/proto"
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"net"
	"os"
	"testing"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/metadata"
)

var configPath = flag.String("config", "../config.yaml", "单元测试配置文件")

func TestMain(m *testing.M) {
	flag.Parse()
	config.InitConf(*configPath)
	mysql.InitMysql()
	//开始测试
	m.Run()
}

// 功能测试以Test开头
func TestAddRecord(t *testing.T) {
	dataList := []*proto.Record{
		{
			Account:         "11111",
			GroupId:         "aaaaa",
			UserMsg:         "你好，golang怎么学",
			UserMsgTokens:   20,
			UserMsgKeywords: []string{"golang", "怎么学"},
			AiMsg:           "对于golang的学习，推荐报名零声教育golang系统课程",
			AiMsgTokens:     50,
			ReqTokens:       100,
			CreateAt:        1683,
			EndpointAccount: "2041993283",
			Endpoint:        0,
			EnterpriseId:    "0voice",
		},
		{
			Account:         "22222",
			GroupId:         "bbbbb",
			UserMsg:         "你好，golang怎么学",
			UserMsgTokens:   20,
			UserMsgKeywords: []string{"golang", "怎么学"},
			AiMsg:           "对于golang的学习，推荐报名零声教育golang系统课程",
			AiMsgTokens:     50,
			ReqTokens:       100,
			CreateAt:        123456789,
			EndpointAccount: "2041993284",
			Endpoint:        1,
			EnterpriseId:    "0voice",
		},
		{
			Account:         "33333",
			GroupId:         "ccccc",
			UserMsg:         "你好，golang怎么学",
			UserMsgTokens:   20,
			UserMsgKeywords: []string{"golang", "怎么学"},
			AiMsg:           "对于golang的学习，推荐报名零声教育golang系统课程",
			AiMsgTokens:     50,
			ReqTokens:       100,
			CreateAt:        123456789,
			EndpointAccount: "2041993285",
			Endpoint:        2,
			EnterpriseId:    "0voice1",
		},
	}

	//启动grpc server
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		t.Error(err)
		return
	}
	//添加服务器启动选项
	var opts = getServerOptions()
	s := grpc.NewServer(opts...)
	defer s.Stop()
	conf := config.GetConf()
	logger := log.NewLogger()
	chatRecordsData := data.NewChatRecordsData(mysql.GetDB(), logger)
	chatGPTDataServer := server.NewChatGPTDataServer(conf, logger, chatRecordsData)
	proto.RegisterChatGPTDataServer(s, chatGPTDataServer)

	go func() {
		err = s.Serve(lis)
		if err != nil {
			t.Error(err)
			return
		}
	}()

	var clientOption = getOptions()

	//客户端调用server
	// conn, err := grpc.Dial("localhost:50051", clientOption...)
	conn, err := grpc.Dial("192.168.10.129:50052", clientOption...)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	client := proto.NewChatGPTDataClient(conn)
	ctx := context.Background()
	ctx = appendBearerTokenToContext(ctx)
	for _, item := range dataList {

		res, err := client.AddRecord(ctx, item)
		if err != nil {
			t.Error(err)
			return
		}
		record, err := chatRecordsData.GetById(res.GetId())
		if err != nil {
			t.Error(err)
			return
		}
		if record.CreateAt != item.CreateAt || record.ReqTokens != int(item.ReqTokens) ||
			record.AIMsgTokens != int(item.AiMsgTokens) || record.UserMsgTokens != int(item.UserMsgTokens) ||
			record.UserMsg != item.UserMsg || record.AIMsg != item.AiMsg ||
			len(record.UserMsgKeywords) != len(item.UserMsgKeywords) || record.GroupID != item.GroupId ||
			record.Account != item.Account {
			t.Error("写入的记录与读取的记录不匹配")
		}

	}
}

func getOptions() []grpc.DialOption {
	opts := make([]grpc.DialOption, 0)
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//opts = append(opts, getTlsOpt("../cert/ca_cert.pem", "chatgpt-data.grpc.0voice.com"))
	//opts = append(opts, getMTlsOpt("../cert/ca_cert.pem", "../cert/client_cert.pem", "../cert/client_key.pem"))
	//opts = append(opts, getAuth())
	return opts
}
func getTlsOpt(cert, serviceName string) grpc.DialOption {
	creds, err := credentials.NewClientTLSFromFile(cert, serviceName)
	if err != nil {
		panic(err)
	}
	return grpc.WithTransportCredentials(creds)
}
func getMTlsOpt(caCert, certFile, keyFile string) grpc.DialOption {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}
	ca := x509.NewCertPool()
	bytes, err := os.ReadFile(caCert)
	if err != nil {
		panic(err)
	}
	ok := ca.AppendCertsFromPEM(bytes)
	if !ok {
		panic("append cert failed")
	}
	tlsConfig := &tls.Config{
		ServerName:   "chatgpt-data.grpc.0voice.com",
		Certificates: []tls.Certificate{cert},
		RootCAs:      ca,
	}
	return grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))
}
func getAuth() grpc.DialOption {
	token := config.GetConf().Server.AccessToken
	perRPC := oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})}
	return grpc.WithPerRPCCredentials(perRPC)
}

func appendBearerTokenToContext(ctx context.Context) context.Context {
	token := config.GetConf().Server.AccessToken
	md := metadata.Pairs("authorization", "Bearer "+token)
	return metadata.NewOutgoingContext(ctx, md)
}

func getServerOptions() []grpc.ServerOption {
	var opts = make([]grpc.ServerOption, 0)
	opts = append(opts, server.GetKeepaliveOpt()...)
	//opts = append(opts, server.GetTlsOpt("../cert/server_cert.pem", "../cert/server_key.pem"))
	//opts = append(opts, server.GetMTlsOpt("../cert/client_ca_cert.pem", "../cert/server_cert.pem", "../cert/server_key.pem"))
	opts = append(opts, grpc.StreamInterceptor(server.StreamInterceptor))
	opts = append(opts, grpc.UnaryInterceptor(server.UnaryInterceptor))
	return opts
}
