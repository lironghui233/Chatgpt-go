package test_case

import (
	chat_context "chatgpt-service/chatgpt-server/chat-context"
	"chatgpt-service/chatgpt-server/server"
	"chatgpt-service/pkg/config"
	"chatgpt-service/pkg/db/redis"
	"chatgpt-service/pkg/log"
	"chatgpt-service/proto"
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var configPath = flag.String("config", "../config.yaml", "单元测试配置文件")

func TestMain(m *testing.M) {
	flag.Parse()
	config.InitConf(*configPath)
	redis.InitRedisPool()

	//启动grpc server
	go startGRPCServer()
	m.Run()
}

func startGRPCServer() {
	cnf := config.GetConf()
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer(getServerOptions()...)
	logger := log.NewLogger()
	logger.SetLevel(cnf.Log.Level)
	logger.SetPrintCaller(true)
	chatGPTServer := server.NewChatGPTServer(cnf, logger)
	proto.RegisterChatGPTServer(s, chatGPTServer)

	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

func getServerOptions() []grpc.ServerOption {
	var opts = make([]grpc.ServerOption, 0)
	opts = append(opts, server.GetKeepaliveOpt()...)
	opts = append(opts, grpc.StreamInterceptor(server.StreamInterceptor))
	opts = append(opts, grpc.UnaryInterceptor(server.UnaryInterceptor))
	return opts
}

func getOptions() []grpc.DialOption {
	opts := make([]grpc.DialOption, 0)
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return opts
}
func appendBearerTokenToContext(ctx context.Context) context.Context {
	token := config.GetConf().Server.AccessToken
	md := metadata.Pairs("authorization", "Bearer "+token)
	return metadata.NewOutgoingContext(ctx, md)

}

func TestChatCompletion(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", getOptions()...)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	client := proto.NewChatGPTClient(conn)
	in := &proto.ChatCompletionRequest{
		Message:       "你好",
		Id:            uuid.New().String(),
		Endpoint:      proto.ChatEndpoint_WEB,
		EnableContext: true,
	}
	ctx := context.Background()
	ctx = appendBearerTokenToContext(ctx)
	res, err := client.ChatCompletion(ctx, in)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res.Choices[0].Message)
	t.Log(res.Usage.TotalTokens)
	t.Log(res.Created)
	fmt.Printf("%+v\n", res)
}
func TestChatCompletionStream(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", getOptions()...)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	client := proto.NewChatGPTClient(conn)
	in := &proto.ChatCompletionRequest{
		Message:       "你好",
		Id:            uuid.New().String(),
		Endpoint:      proto.ChatEndpoint_WEB,
		EnableContext: true,
	}
	ctx := context.Background()
	ctx = appendBearerTokenToContext(ctx)
	stream, err := client.ChatCompletionStream(ctx, in)
	if err != nil {
		t.Error(err)
		return
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error(err)
			return
		}
		// t.Log(res.Id)
		// t.Log(res.Created)
		// t.Log(res.Choices[0].Delta)
		fmt.Printf("%+v\n", res)
	}
	stream.CloseSend()
}

func TestChatCompletionQQ(t *testing.T) {
	enterpriseId := "leoh_enterprise"
	endpointAccount := "2041993283"
	id := "13579101112"
	group := "qwertyu"
	dataList := []string{
		"你好",
		"hello",
		"k8s 是什么",
	}
	conn, err := grpc.Dial("localhost:50051", getOptions()...)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	client := proto.NewChatGPTClient(conn)
	for _, msg := range dataList {
		in := &proto.ChatCompletionRequest{
			Message:         msg,
			Id:              id,
			GroupId:         group,
			Endpoint:        proto.ChatEndpoint_QQ,
			EnableContext:   true,
			EndpointAccount: endpointAccount,
			EnterpriseId:    enterpriseId,
		}
		ctx := context.Background()
		ctx = appendBearerTokenToContext(ctx)
		_, err := client.ChatCompletion(ctx, in)
		if err != nil {
			t.Error(err)
			return
		}
	}
	time.Sleep(time.Second)
	cache := chat_context.GetCacheContext(proto.ChatEndpoint_QQ)
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
	if len(cmList) > 2 {
		if dataList[len(dataList)-1] != cmList[1].Message.Content {
			t.Error("上下文获取有误")
		}
	}
	for _, item := range cmList {
		t.Log(item.TokensNum, item.Message.Role, item.Message.Content)
	}
}

func TestChatCompletionStreamWeb(t *testing.T) {
	enterpriseId := "leoh_enterprise"
	endpointAccount := "2041993283"
	dataList := []string{
		"你好",
		"hello",
		"docker 是什么",
	}
	conn, err := grpc.Dial("localhost:50051", getOptions()...)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	pid := ""
	client := proto.NewChatGPTClient(conn)
	for _, msg := range dataList {
		in := &proto.ChatCompletionRequest{
			Message:         msg,
			Id:              uuid.New().String(),
			GroupId:         "",
			Endpoint:        proto.ChatEndpoint_WEB,
			EnableContext:   true,
			Pid:             pid,
			EndpointAccount: endpointAccount,
			EnterpriseId:    enterpriseId,
		}
		ctx := context.Background()
		ctx = appendBearerTokenToContext(ctx)
		stream, err := client.ChatCompletionStream(ctx, in)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Error(err)
				return
			}
			pid = res.Id
		}
		stream.CloseSend()
	}
	time.Sleep(time.Second)
	cache := chat_context.GetCacheContext(proto.ChatEndpoint_WEB)
	l, err := cache.Get(pid, "", proto.ChatEndpoint_WEB)
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
	if len(cmList) > 2 {
		if dataList[len(dataList)-1] != cmList[1].Message.Content {
			t.Error("上下文获取有误")
		}
	}
	for _, item := range cmList {
		t.Log(item.TokensNum, item.Message.Role, item.Message.Content)
	}
}
