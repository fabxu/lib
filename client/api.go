package client

import (
	"github.com/fabxu/lib/client/http"
	"github.com/fabxu/lib/client/redis"
	"github.com/fabxu/lib/client/sqldb"
)

var (
	SQLDB = &sqldb.Client{}
	Redis = &redis.Client{}
	HTTP  = &http.Wrapper{}
)
