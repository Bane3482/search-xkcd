package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	searchpb "yadro.com/course/proto/search"
	"yadro.com/course/search/core"
)

type Server struct {
	searchpb.UnimplementedSearchServer
	service core.Searcher
}

func NewServer(service core.Searcher) *Server {
	return &Server{service: service}
}

func (s *Server) Ping(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *Server) Search(ctx context.Context, req *searchpb.SearchRequest) (*searchpb.SearchReply, error) {
	if req.Phrase == "" {
		return nil, status.Error(codes.InvalidArgument, "empty phrase")
	}

	reply, err := s.service.Search(ctx, req.Phrase, int(req.Limit))

	if err != nil {
		return nil, err
	}

	return &searchpb.SearchReply{
		Comics: core.ToProtoComics(reply),
		Total:  int64(len(reply)),
	}, nil
}

func (s *Server) ISearch(ctx context.Context, req *searchpb.SearchRequest) (*searchpb.SearchReply, error) {
	if req.Phrase == "" {
		return nil, status.Error(codes.InvalidArgument, "empty phrase")
	}

	reply, err := s.service.ISearch(ctx, req.Phrase, int(req.Limit))

	if err != nil {
		return nil, err
	}

	return &searchpb.SearchReply{
		Comics: core.ToProtoComics(reply),
		Total:  int64(len(reply)),
	}, nil
}
