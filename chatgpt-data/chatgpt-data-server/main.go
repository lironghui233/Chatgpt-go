package main

import (
	"chatgpt-data/chatgpt-data-server/data"
	"chatgpt-data/chatgpt-data-server/server"
	"chatgpt-data/pkg/cmd"
	"chatgpt-data/pkg/config"
	"chatgpt-data/pkg/db/mysql"
	"chatgpt-data/pkg/log"
	"chatgpt-data/proto"
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

	// s := grpc.NewServer(server.GetOptions()...)
	s := grpc.NewServer()

	logger := log.NewLogger()
	logger.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	logger.SetLevel(cnf.Log.Level)
	logger.SetPrintCaller(true)

	chatRecordsData := data.NewChatRecordsData(mysql.GetDB(), logger)
	proto.RegisterChatGPTDataServer(s, server.NewChatGPTDataServer(cnf, logger, chatRecordsData))

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
	log.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	log.SetLevel(cnf.Log.Level)
	mysql.InitMysql()
}
