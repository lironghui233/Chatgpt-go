package client

import (
	"chatgpt-wechat/pkg/log"
	"chatgpt-wechat/services"
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

type ServiceClient interface {
	GetPool(addr string) services.ClientPool
}

type DefaultClient struct {
}

func AppendBearerTokenToContext(ctx context.Context, accessToken string) context.Context {
	md := metadata.Pairs("authorization", "Bearer "+accessToken)
	return metadata.NewOutgoingContext(ctx, md)
}

func (c *DefaultClient) getOptions() []grpc.DialOption {
	opts := make([]grpc.DialOption, 0)
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, c.getKeepaliveOpt())
	return opts
}

func (c *DefaultClient) getKeepaliveOpt() (opt grpc.DialOption) {
	var kacp = keepalive.ClientParameters{
		// 如果没有活动（请求，流），则每30s发送一次ping
		Time: 30 * time.Second,
		// ping ack 1s 内没有返回这认定为连接断开
		Timeout: time.Second,
		// 当没有任何活动流的情况下，是否允许被ping
		PermitWithoutStream: true,
	}
	return grpc.WithKeepaliveParams(kacp)
}

func (c *DefaultClient) GetPool(addr string) services.ClientPool {
	pool, err := services.GetPool(addr, c.getOptions()...)
	if err != nil {
		log.Error(err)
		return nil
	}
	return pool
}
