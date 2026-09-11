package xkcd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"yadro.com/course/update/core"
)

type Client struct {
	log    *slog.Logger
	client http.Client
	url    string
}

const maxRetries = 5

const infoPath = "/info.0.json"

func NewClient(url string, timeout time.Duration, log *slog.Logger) (*Client, error) {
	if url == "" {
		return nil, fmt.Errorf("empty base url specified")
	}
	return &Client{
		client: http.Client{Timeout: timeout},
		log:    log,
		url:    url,
	}, nil
}

func (c Client) Get(ctx context.Context, id int) (core.XKCDInfo, error) {
	return c.getByUrl(ctx, fmt.Sprintf("%s/%d/%s", c.url, id, infoPath))
}

func (c Client) LastID(ctx context.Context) (int, error) {
	result, err := c.getByUrl(ctx, c.url+infoPath)

	if err != nil {
		c.log.Error("xkcd client last id get", "error", err)
		return 0, err
	}

	return result.ID, nil
}

func (c Client) getByUrl(ctx context.Context, url string) (core.XKCDInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return core.XKCDInfo{}, err
	}

	var resp *http.Response

	for range maxRetries {
		resp, err = c.client.Do(req)

		if err == nil {
			break
		}
	}

	if err != nil {
		c.log.Error("xkcd get by url request", "error", err)
		return core.XKCDInfo{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		c.log.Error("xkcd get by url status", "error", err)
		return core.XKCDInfo{}, core.ErrNotFound
	}

	var result core.XKCDInfo

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.log.Error("xkcd get by url decode", "error", err)
		return core.XKCDInfo{}, err
	}

	return result, nil
}
