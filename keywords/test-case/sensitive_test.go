package test_case

import (
	"context"
	"flag"
	"fmt"
	"keywords/keywords-server/server"
	"keywords/pkg/config"
	"keywords/pkg/filter"
	"keywords/proto"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var configPath = flag.String("config", "../config.yaml", "单元测试配置文件")
var dictPath = flag.String("dict", "dict.txt", "单元测试关键词库文件")

func TestMain(m *testing.M) {
	flag.Parse()
	config.InitConf(*configPath)
	filter.InitFilter(*dictPath)
	m.Run()
}
func TestFindAll(t *testing.T) {
	dataList := []struct {
		req *proto.FindAllReq
		res *proto.FindAllRes
	}{
		{
			req: &proto.FindAllReq{
				Text: "abcdefg defer recover sync Protobuf abcdefg",
			},
			res: &proto.FindAllRes{
				Keywords: []string{"defer", "recover", "sync", "Protobuf"},
			},
		},
		{
			req: &proto.FindAllReq{
				Text: "defer recover sync Protobuf",
			},
			res: &proto.FindAllRes{
				Keywords: []string{"defer", "recover", "sync", "Protobuf"},
			},
		},
		{
			req: &proto.FindAllReq{
				Text: "sync论文猎狗Protobuf",
			},
			res: &proto.FindAllRes{
				Keywords: []string{"sync", "论文", "猎狗"},
			},
		},
		{
			req: &proto.FindAllReq{
				Text: "你可能需要一篇golang论文",
			},
			res: &proto.FindAllRes{
				Keywords: []string{"golang论文", "论文"},
			},
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
	filter := filter.GetFilter()
	proto.RegisterKeywordsServer(s, server.NewKeywordsServer(filter))
	go func() {
		err = s.Serve(lis)
		if err != nil {
			t.Error(err)
			return
		}
	}()
	var clientOption = getOptions()
	//客户端调用server
	conn, err := grpc.Dial("localhost:50051", clientOption...)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	client := proto.NewKeywordsClient(conn)
	ctx := context.Background()
	ctx = appendBearerTokenToContext(ctx)
	for _, item := range dataList {
		res, err := client.FindAll(ctx, item.req)
		if err != nil {
			t.Error(err)
		}
		if len(res.Keywords) != len(item.res.Keywords) {
			fmt.Println(res.Keywords)
			fmt.Println(item.res.Keywords)
			t.Error("关键词过滤结果与预期不一致")
		}
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
