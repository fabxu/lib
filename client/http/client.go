package http

import (
	"context"
	"time"

	"github.com/go-resty/resty/v2"
)

func NewClient(_ context.Context, cfg Config) *resty.Client {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 10
	}

	return resty.New().
		SetBaseURL(cfg.Host).
		SetTimeout(time.Duration(timeout) * time.Second)
}
