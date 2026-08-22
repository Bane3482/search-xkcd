package search

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"yadro.com/course/api/core"
	searchpb "yadro.com/course/proto/search"
)

type Client struct {
	log    *slog.Logger
	client searchpb.SearchClient
}

func New(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Error("search client connection", "error", err)
	}

	return &Client{
		log:    log,
		client: searchpb.NewSearchClient(conn),
	}, nil
}

func (c Client) Ping(ctx context.Context) error {
	_, err := c.client.Ping(ctx, nil)

	if err != nil {
		c.log.Error("searh client ping", "error", err)
		return err
	}

	return nil
}

func (c Client) Search(ctx context.Context, phrase string, limit int) (core.ComicsReply, error) {
	reply, err := c.client.Search(ctx, &searchpb.SearchRequest{Phrase: phrase, Limit: int64(limit)})

	if err != nil {
		c.log.Error("search client search", "error", err)
		return core.ComicsReply{}, err
	}

	return core.FromProtoComicsReply(reply), nil
}

func (c Client) ISearch(ctx context.Context, phrase string, limit int) (core.ComicsReply, error) {
	reply, err := c.client.ISearch(ctx, &searchpb.SearchRequest{Phrase: phrase, Limit: int64(limit)})

	if err != nil {
		c.log.Error("search client search", "error", err)
		return core.ComicsReply{}, err
	}

	return core.FromProtoComicsReply(reply), nil
}
