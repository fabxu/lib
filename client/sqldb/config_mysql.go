package sqldb

import (
	"fmt"
	"os"
	"strings"

	"github.com/fabxu/lib/client/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLConfig struct {
	Config `mapstructure:",squash"`
}

func (cfg MySQLConfig) GetDSN(useDB bool) (string, error) {
	addr := cfg.Addr
	addrEnv := "MYSQL_ADDRESS"

	if cfg.AddressEnv != "" {
		addrEnv = cfg.AddressEnv
	}

	if v, ok := os.LookupEnv(addrEnv); ok {
		addr = v
	}

	dbName := cfg.DBName
	databaseEnv := "MYSQL_DATABASE"

	if cfg.DatabaseEnv != "" {
		databaseEnv = cfg.DatabaseEnv
	}

	if v, ok := os.LookupEnv(databaseEnv); ok {
		dbName = v
	}

	username := cfg.Username
	usernameEnv := "MYSQL_USERNAME"

	if cfg.UsernameEnv != "" {
		usernameEnv = cfg.UsernameEnv
	}

	if v, ok := os.LookupEnv(usernameEnv); ok {
		username = v
	}

	password := cfg.Password
	passwordEnv := "MYSQL_PASSWORD"

	if cfg.PasswordEnv != "" {
		passwordEnv = cfg.PasswordEnv
	}

	if v, ok := os.LookupEnv(passwordEnv); ok {
		password = v
	}

	protocol := strings.TrimSpace(cfg.Protocol)
	if protocol == "" {
		protocol = "tcp"
	}

	conf := util.ConfigMap{
		"MYSQL_USERNAME": username,
		"MYSQL_PASSWORD": password,
		"MYSQL_ADDRESS":  addr,
		"MYSQL_DATABASE": dbName,
	}.TrimSpace()
	if has, key := conf.AnyEmpty(); has {
		return "", fmt.Errorf("missing mysql necessary param %q, please check your config and try again", key)
	}

	param := ""
	if p := strings.TrimSpace(cfg.Param); p != "" {
		param = fmt.Sprintf("?%v", p)
	}

	if !useDB {
		return fmt.Sprintf("%v:%v@%v(%v)/%v", username, password, protocol, addr, param), nil
	}

	return fmt.Sprintf("%v:%v@%v(%v)/%v%v", username, password, protocol, addr, dbName, param), nil
}

func (cfg MySQLConfig) GetDialector() (func(string) gorm.Dialector, error) {
	return mysql.Open, nil
}
