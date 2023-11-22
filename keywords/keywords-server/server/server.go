package server

import (
	"context"
	"keywords/pkg/filter"
	"keywords/proto"
)

type keywordsServer struct {
	proto.UnimplementedKeywordsServer
	filter filter.IKeywordsFilter
}

func NewKeywordsServer(filter filter.IKeywordsFilter) proto.KeywordsServer {
	return &keywordsServer{
		filter: filter,
	}
}
func (s *keywordsServer) FindAll(ctx context.Context, in *proto.FindAllReq) (*proto.FindAllRes, error) {
	list := s.filter.FindAll(in.Text)
	return &proto.FindAllRes{
		Keywords: list,
	}, nil
}
