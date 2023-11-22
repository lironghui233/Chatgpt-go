package server

import (
	"context"
	"sensitive-words/pkg/filter"
	"sensitive-words/proto"
)

type SensitiveWordsServer struct {
	proto.UnimplementedSensitiveWordsServer
	filter filter.ISensitiveFilter
}

func NewSensitiveWordsServer(filter filter.ISensitiveFilter) proto.SensitiveWordsServer {
	return &SensitiveWordsServer{
		filter: filter,
	}
}
func (s *SensitiveWordsServer) Validate(ctx context.Context, in *proto.ValidateReq) (*proto.ValidateRes, error) {
	ok, word := s.filter.Validate(in.Text)
	return &proto.ValidateRes{
		Ok:   ok,
		Word: word,
	}, nil
}
