package sqldb

import (
	"fmt"

	"gorm.io/gorm"
)

type DBType string

const (
	DBTypeMySQL    DBType = "mysql"
	DBTypePostgres DBType = "postgres"
)

type ConfigGetter interface {
	GetDSN(useDB bool) (string, error)
	GetDialector() (func(string) gorm.Dialector, error)
}

type Config struct {
	DBType       DBType
	Username     string
	Password     string
	Protocol     string
	Addr         string
	DBName       string
	Param        string
	MaxOpenConns int  `mapstructure:"max_open_conns"`
	MaxIdleConns int  `mapstructure:"max_idle_conns"`
	MaxIdleTime  int  `mapstructure:"max_idle_time"`
	ShowStats    bool `mapstructure:"show_stats"`
	CreateDB     bool `mapstructure:"create_db"`
	Verbose      bool
	AddressEnv   string
	DatabaseEnv  string
	UsernameEnv  string
	PasswordEnv  string
}

func (cfg Config) GetDSN(useDB bool) (string, error) {
	switch cfg.DBType {
	case DBTypeMySQL:
		return cfg.ToMySQL().GetDSN(useDB)
	case DBTypePostgres:
		return cfg.ToPostgres().GetDSN(useDB)
	default:
		return "", fmt.Errorf("unsupported db type")
	}
}

func (cfg Config) GetDialector() (func(string) gorm.Dialector, error) {
	switch cfg.DBType {
	case DBTypeMySQL:
		return cfg.ToMySQL().GetDialector()
	case DBTypePostgres:
		return cfg.ToPostgres().GetDialector()
	default:
		return nil, fmt.Errorf("unsupported dialector")
	}
}

func (cfg Config) ToMySQL() MySQLConfig {
	return MySQLConfig{Config: cfg}
}

func (cfg Config) ToPostgres() PostgresConfig {
	return PostgresConfig{Config: cfg}
}

func (cfg Config) chooseSourcePath(sourcePaths SourcePathMapper) (string, error) {
	srcPath, ok := sourcePaths[cfg.DBType]
	if !ok {
		return "", fmt.Errorf("no define source path")
	}

	return srcPath, nil
}

type SourcePathMapper map[DBType]string

func NewSourcePathMapper() SourcePathMapper {
	return make(map[DBType]string)
}

func (spm SourcePathMapper) ForMySQL(sourcePath string) SourcePathMapper {
	spm[DBTypeMySQL] = sourcePath
	return spm
}

func (spm SourcePathMapper) ForPostgres(sourcePath string) SourcePathMapper {
	spm[DBTypePostgres] = sourcePath
	return spm
}
