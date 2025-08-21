package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	*redis.Client
}

func New(ctx context.Context, cfg Config) *Client {
	return &Client{
		Client: NewDB(ctx, cfg),
	}
}

func (c *Client) Global(ctx context.Context, cfg Config) {
	// 确保receiver的引用内容变更
	*c = *New(ctx, cfg)
}
