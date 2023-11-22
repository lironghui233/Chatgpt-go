package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"sensitive-words/pkg/cmd"
	"sensitive-words/pkg/config"
	"sensitive-words/pkg/filter"
	"sensitive-words/proto"
	"sensitive-words/sensitive-server/server"
)

func main() {
	loadDependOn()
	if cmd.Args.InitDict {
		filter.OverwriteDict(cmd.Args.Dict)
		return
	}
	cnf := config.GetConf()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cnf.Server.Host, cnf.Server.Port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(server.GetOptions()...)
	proto.RegisterSensitiveWordsServer(s, server.NewSensitiveWordsServer(filter.GetFilter()))

	//添加健康检查逻辑
	healthCheck := health.NewServer()
	grpchealth.RegisterHealthServer(s, healthCheck)

	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

func loadDependOn() {
	config.InitConf(cmd.Args.Config)
	filter.InitFilter(cmd.Args.Dict)
}
