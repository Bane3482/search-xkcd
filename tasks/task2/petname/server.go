package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"yadro.com/course/config"
	petnamepb "yadro.com/course/proto"
)

type generator interface {
	Generate(words int64, separator string) string
	GenerateMany(words int64, separator string, names int64) []string
}

type server struct {
	generator generator
	petnamepb.UnimplementedPetnameGeneratorServer
}

func (s *server) Ping(_ context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *server) Generate(ctx context.Context, req *petnamepb.PetnameRequest) (*petnamepb.PetnameResponse, error) {
	params := FromProtoParams(req)

	if params.Words <= 0 {
		return nil, status.Error(codes.InvalidArgument, "words invalid")
	}

	return &petnamepb.PetnameResponse{
		Name: s.generator.Generate(params.Words, params.Separator),
	}, nil
}

func (s *server) GenerateMany(req *petnamepb.PetnameStreamRequest, server grpc.ServerStreamingServer[petnamepb.PetnameResponse]) error {
	params := FromProtoStream(req)

	if params.Names <= 0 {
		return status.Error(codes.InvalidArgument, "names invalid")
	}

	if params.WordParams.Words <= 0 {
		return status.Error(codes.InvalidArgument, "words invalid")
	}

	result := s.generator.GenerateMany(params.WordParams.Words, params.WordParams.Separator, params.Names)

	for _, name := range result {
		resp := &petnamepb.PetnameResponse{
			Name: name,
		}

		if err := server.Send(resp); err != nil {
			return status.Error(codes.Internal, "server send error")
		}
	}

	return nil
}

func main() {
	configPath := flag.String("config", "", "config for setting app")

	cfg, err := config.New(*configPath)

	if err != nil {
		fmt.Println("error config")
		return
	}

	listener, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	petnamepb.RegisterPetnameGeneratorServer(s, &server{})
	reflection.Register(s)

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type WordParams struct {
	Words     int64
	Separator string
}

type WordStreamParams struct {
	Names      int64
	WordParams WordParams
}

func FromProtoParams(e *petnamepb.PetnameRequest) *WordParams {
	return &WordParams{
		Words:     e.Words,
		Separator: e.Separator,
	}
}

func FromProtoStream(e *petnamepb.PetnameStreamRequest) *WordStreamParams {
	return &WordStreamParams{
		Names: e.Names,
		WordParams: WordParams{
			Words:     e.Words,
			Separator: e.Separator,
		},
	}
}
