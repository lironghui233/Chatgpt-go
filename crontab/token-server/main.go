package main

import (
	"crontab/cron"
	"crontab/pkg/cmd"
	"crontab/pkg/config"
	"crontab/pkg/db/redis"
	"crontab/pkg/log"
	"crontab/proto"
	server "crontab/token-server/server"
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
		log.Error(err)
		panic(err)
	}
	s := grpc.NewServer(server.GetOptions()...)

	logger := log.NewLogger()
	logger.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	logger.SetLevel(cnf.Log.Level)
	logger.SetPrintCaller(true)

	tokenServer := server.NewTokenServer(cnf, logger)
	proto.RegisterTokenServer(s, tokenServer)

	//添加健康检查
	healthCheck := health.NewServer()
	grpchealth.RegisterHealthServer(s, healthCheck)

	//启动自动任务，刷新token
	go cron.Run()

	if err := s.Serve(lis); err != nil {
		panic(err)
	}

}

func loadDependOn() {
	//初始化配置
	config.InitConf(cmd.Args.Config)
	cnf := config.GetConf()

	//初始化日志
	log.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	log.SetLevel(cnf.Log.Level)

	//初始化redis
	redis.InitRedisPool()

}
