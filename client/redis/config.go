package redis

import (
	"fmt"
	"github.com/fabxu/lib/client/util"
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
	"strings"
)

const (
	TypeSingle   = "single"
	TypeSentinel = "sentinel"
)

type Config struct {
	Type   string
	Single struct {
		Addr     string
		Username string
		Password string
		DB       int
	}
	Sentinel struct {
		Master   string
		Addrs    []string
		Username string
		Password string
		DB       int
	}
}

func (cfg Config) GetSingleOptions() (*redis.Options, error) {
	addr := cfg.Single.Addr
	if v, ok := os.LookupEnv("REDIS_ADDRESS"); ok {
		addr = v
	}

	username := cfg.Single.Username
	if v, ok := os.LookupEnv("REDIS_USERNAME"); ok {
		username = v
	}

	password := cfg.Single.Password
	if v, ok := os.LookupEnv("REDIS_PASSWORD"); ok {
		password = v
	}

	db := cfg.Single.DB
	if v, ok := os.LookupEnv("REDIS_DATABASE"); ok {
		db, _ = strconv.Atoi(v)
	}

	conf := util.ConfigMap{
		"REDIS_ADDRESS":  addr,
		"REDIS_DATABASE": strconv.Itoa(db),
	}.TrimSpace()
	if has, key := conf.AnyEmpty(); has {
		return nil, fmt.Errorf("missing redis necessary param %q, please check your config and try again", key)
	}

	return &redis.Options{Addr: addr, Username: username, Password: password, DB: db}, nil
}

func (cfg Config) GetSentinelOptions() (*redis.FailoverOptions, error) {
	master := cfg.Sentinel.Master
	if v, ok := os.LookupEnv("REDIS_MASTER"); ok {
		master = v
	}

	addrs := cfg.Sentinel.Addrs
	if v, ok := os.LookupEnv("REDIS_ADDRESS"); ok {
		addrs = strings.Split(v, ";")
	}

	username := cfg.Sentinel.Username
	if v, ok := os.LookupEnv("REDIS_USERNAME"); ok {
		username = v
	}

	password := cfg.Sentinel.Password
	if v, ok := os.LookupEnv("REDIS_PASSWORD"); ok {
		password = v
	}

	db := cfg.Sentinel.DB
	if v, ok := os.LookupEnv("REDIS_DATABASE"); ok {
		db, _ = strconv.Atoi(v)
	}

	conf := util.ConfigMap{
		"REDIS_MASTER":   master,
		"REDIS_ADDRESS":  strings.Join(addrs, ";"),
		"REDIS_DATABASE": strconv.Itoa(db),
	}.TrimSpace()
	if has, key := conf.AnyEmpty(); has {
		return nil, fmt.Errorf("missing redis necessary param %q, please check your config and try again", key)
	}

	return &redis.FailoverOptions{
		MasterName: master, SentinelAddrs: addrs,
		Username: username, Password: password, DB: db,
	}, nil
}
