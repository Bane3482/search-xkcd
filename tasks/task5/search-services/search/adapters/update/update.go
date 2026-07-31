package update

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	updatepb "yadro.com/course/proto/update"
	"yadro.com/course/search/core"
)

type Client struct {
	log    *slog.Logger
	client updatepb.UpdateClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Error("client connection", "error", err)
	}

	return &Client{
		log:    log,
		client: updatepb.NewUpdateClient(conn),
	}, nil
}

func (c Client) Get(ctx context.Context) ([]*core.Comics, error) {
	reply, err := c.client.Get(ctx, nil)

	if err != nil {
		c.log.Error("update client get", "error", err)
		return nil, err
	}

	return core.FromProtoComics(reply.Comics), nil
}
