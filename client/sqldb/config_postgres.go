package sqldb

import (
	"fmt"
	"os"
	"strings"

	"github.com/fabxu/lib/client/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Config `mapstructure:",squash"`
}

func (cfg PostgresConfig) GetDSN(useDB bool) (string, error) {
	addr := cfg.Addr
	if v, ok := os.LookupEnv("POSTGRES_ADDRESS"); ok {
		addr = v
	}

	dbName := cfg.DBName
	if v, ok := os.LookupEnv("POSTGRES_DATABASE"); ok {
		dbName = v
	}

	username := cfg.Username
	if v, ok := os.LookupEnv("POSTGRES_USERNAME"); ok {
		username = v
	}

	password := cfg.Password
	if v, ok := os.LookupEnv("POSTGRES_PASSWORD"); ok {
		password = v
	}

	addrs := strings.Split(addr, ":")
	if len(addrs) < 2 {
		return "", fmt.Errorf("invalid postgres address")
	}

	host := addrs[0]
	port := addrs[1]

	conf := util.ConfigMap{
		"POSTGRES_HOST":     host,
		"POSTGRES_PORT":     port,
		"POSTGRES_USERNAME": username,
		"POSTGRES_PASSWORD": password,
		"POSTGRES_DATABASE": dbName,
	}.TrimSpace()
	if has, key := conf.AnyEmpty(); has {
		return "", fmt.Errorf("missing postgres necessary param %q, please check your config and try again", key)
	}

	param := ""
	if p := strings.TrimSpace(cfg.Param); p != "" {
		param = strings.ReplaceAll(p, " ", "&")
	}

	if !useDB {
		return fmt.Sprintf("postgres://%v:%v@%v?%v", username, password, addr, param), nil
	}

	return fmt.Sprintf("postgres://%v:%v@%v/%v?%v", username, password, addr, dbName, param), nil
}

func (cfg PostgresConfig) GetDialector() (func(string) gorm.Dialector, error) {
	return postgres.Open, nil
}
