package main

import (
	"chatgpt-service/chatgpt-server/server"
	"chatgpt-service/pkg/cmd"
	"chatgpt-service/pkg/config"
	"chatgpt-service/pkg/db/redis"
	"chatgpt-service/pkg/log"
	"chatgpt-service/proto"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	loadDependOn()
	cnf := config.GetConf()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cnf.Server.Host, cnf.Server.Port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(server.GetOptions()...)

	logger := log.NewLogger()
	logger.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	logger.SetLevel(cnf.Log.Level)
	logger.SetPrintCaller(true)
	chatGPTServer := server.NewChatGPTServer(cnf, logger)
	proto.RegisterChatGPTServer(s, chatGPTServer)

	//添加健康检查逻辑
	healthCheck := health.NewServer()
	grpchealth.RegisterHealthServer(s, healthCheck)

	if err := s.Serve(lis); err != nil {
		panic(err)
	}

}
func loadDependOn() {
	config.InitConf(cmd.Args.Config)
	cnf := config.GetConf()
	redis.InitRedisPool()
	log.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	log.SetLevel(cnf.Log.Level)

}
