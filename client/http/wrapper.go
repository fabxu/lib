package http

import (
	"context"
	cmlog "github.com/fabxu/log"
	"github.com/go-resty/resty/v2"
)

type Wrapper struct {
	clients map[string]*resty.Client
	cfgs    map[string]Config
}

func New(ctx context.Context, key string, cfg Config) *Wrapper {
	return &Wrapper{
		clients: map[string]*resty.Client{
			key: NewClient(ctx, cfg),
		},
		cfgs: map[string]Config{
			key: cfg,
		},
	}
}

func (c *Wrapper) Global(ctx context.Context, cfgs map[string]Config) {
	logger := cmlog.Extract(ctx)

	c.cfgs = cfgs
	c.clients = map[string]*resty.Client{}

	for key, cfg := range cfgs {
		logger.Debugf("initializing %v http client", key)
		c.clients[key] = New(ctx, key, cfg).GetClient(key)
	}
}

func (c *Wrapper) GetClient(key string) *resty.Client {
	return c.clients[key]
}
