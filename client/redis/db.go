package redis

import (
	"context"
	cmlog "github.com/fabxu/log"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

func NewDB(ctx context.Context, cfg Config) *redis.Client {
	logger := cmlog.Extract(ctx)

	var rdb *redis.Client

	retry := 0

	for {
		if retry > 10 {
			logger.Panic("redis connect max retry exceeded")
		}

		if strings.ToLower(cfg.Type) == TypeSingle {
			opts, err := cfg.GetSingleOptions()
			if err == nil {
				rdb = redis.NewClient(opts)

				err = rdb.Ping(context.Background()).Err()
				if err == nil {
					break
				}
			}

			logger.Errorf("connecting to redis failed: %v", err)
		} else if strings.ToLower(cfg.Type) == TypeSentinel {
			opts, err := cfg.GetSentinelOptions()
			if err == nil {
				rdb = redis.NewFailoverClient(opts)

				err = rdb.Ping(context.Background()).Err()
				if err == nil {
					break
				}
			}

			logger.Errorf("connecting to redis failed: %v", err)
		}

		logger.Infof("redis connect retry [%v]th", retry+1)
		time.Sleep(5 * time.Second)

		retry++
	}

	logger.Info("connected to redis")

	return rdb
}
