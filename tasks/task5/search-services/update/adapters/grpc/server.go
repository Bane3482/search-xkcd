package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	updatepb "yadro.com/course/proto/update"
	"yadro.com/course/update/core"
)

func NewServer(service core.Updater) *Server {
	return &Server{service: service}
}

type Server struct {
	updatepb.UnimplementedUpdateServer
	service core.Updater
}

func (s *Server) Ping(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *Server) Status(ctx context.Context, req *emptypb.Empty) (*updatepb.StatusReply, error) {
	status := s.service.Status(ctx)

	return &updatepb.StatusReply{
		Status: core.ToProtoStatus(status),
	}, nil
}

func (s *Server) Update(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	if err := s.service.Update(ctx); err != nil {
		if errors.Is(err, core.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "update task already running")
		}
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) Stats(ctx context.Context, req *emptypb.Empty) (*updatepb.StatsReply, error) {
	stats, err := s.service.Stats(ctx)

	if err != nil {
		return nil, err
	}

	return core.ToProtoStats(stats), nil
}

func (s *Server) Drop(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if err := s.service.Drop(ctx); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
