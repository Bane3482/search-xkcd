package update

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"yadro.com/course/api/core"
	updatepb "yadro.com/course/proto/update"
)

type Client struct {
	log    *slog.Logger
	client updatepb.UpdateClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		client: updatepb.NewUpdateClient(conn),
		log:    log,
	}, nil
}

func (c Client) Ping(ctx context.Context) error {
	_, err := c.client.Ping(ctx, nil)

	if err != nil {
		c.log.Error("update client ping")
	}

	return err
}

func (c Client) Status(ctx context.Context) (core.StatusUpdate, error) {
	resp, err := c.client.Status(ctx, nil)

	if err != nil {
		c.log.Error("update client status", "error", err)
		return core.StatusUpdateUnknown, err
	}

	return core.FromProtoStatus(resp.Status), nil
}

func (c Client) Stats(ctx context.Context) (core.UpdateStats, error) {
	resp, err := c.client.Stats(ctx, nil)

	if err != nil {
		c.log.Error("update client stats", "error", err)
		return core.UpdateStats{}, err
	}

	return core.UpdateStats{
		WordsTotal:    int(resp.WordsTotal),
		WordsUnique:   int(resp.WordsUnique),
		ComicsFetched: int(resp.ComicsFetched),
		ComicsTotal:   int(resp.ComicsTotal),
	}, nil

}

func (c Client) Update(ctx context.Context) error {
	_, err := c.client.Update(ctx, nil)

	if err != nil {
		c.log.Error("update client update", "error", err)
		if status.Code(err) == codes.AlreadyExists {
			return core.ErrAlreadyExists
		}
	}
	return err
}

func (c Client) Drop(ctx context.Context) error {
	_, err := c.client.Drop(ctx, nil)

	if err != nil {
		c.log.Error("update client drop", "error", err)
	}

	return err
}
