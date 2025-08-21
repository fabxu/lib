package sqldb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	cmlog "github.com/fabxu/log"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	// initialize db
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Client struct {
	*gorm.DB
}

func New(ctx context.Context, cfg Config) *Client {
	client := &Client{}
	db := client.openDB(ctx, cfg, true)
	client.DB = db

	return client
}

func (c *Client) Global(ctx context.Context, cfg Config) {
	// 确保receiver的引用内容变更
	*c = *New(ctx, cfg)
}

func (c *Client) AutoMigrate(dst ...interface{}) error {
	if CheckDBType(c.DB) == DBTypeMySQL {
		return c.DB.Set(
			"gorm:table_options",
			"ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci",
		).AutoMigrate(dst...)
	}

	return c.DB.AutoMigrate(dst...)
}

func (c *Client) MigrateUp(ctx context.Context, cfg Config, sourcePaths SourcePathMapper) error {
	sourcePath, err := cfg.chooseSourcePath(sourcePaths)
	if err != nil {
		return err
	}

	return c.migrateUp(ctx, cfg, sourcePath)
}

func (c *Client) Transaction(ctx context.Context, f func(ctx context.Context) (interface{}, error)) (res interface{}, err error) {
	tx := c.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()

			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
		}
	}()

	err = tx.Error
	if err != nil {
		return nil, err
	}

	res, err = f(c.Inject(ctx, tx))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return res, nil
}

func (c *Client) openDB(ctx context.Context, cfg Config, useDB bool) *gorm.DB {
	logger := cmlog.Extract(ctx)

	dsn := c.getDSN(ctx, cfg, useDB)
	dbOpener := c.getDBOpener(ctx, cfg)
	retry := 0

	dbLogger := gormlog.Default
	if cfg.Verbose {
		dbLogger = dbLogger.LogMode(gormlog.Info)
	}

	var (
		db  *gorm.DB
		err error
	)

	for {
		if retry > 10 {
			logger.Panicf("%s connect max retry exceeded", cfg.DBType)
		}

		db, err = gorm.Open(dbOpener(dsn), &gorm.Config{Logger: dbLogger})
		if err == nil {
			break
		}

		logger.Errorf("connecting to %s [%s - %s] failed: %v", cfg.DBType, cfg.Addr, cfg.DBName, err)
		logger.Infof("%s connect retry [%v]th", cfg.DBType, retry+1)
		time.Sleep(5 * time.Second)

		retry++
	}

	maxOpenConns, maxIdleConns := 50, 25
	maxIdleTime := 3600

	if cfg.MaxOpenConns > 0 {
		maxOpenConns = cfg.MaxOpenConns
	}

	if cfg.MaxIdleConns > 0 {
		maxIdleConns = cfg.MaxIdleConns
	}

	if cfg.MaxIdleConns > 0 {
		maxIdleTime = cfg.MaxIdleConns
	}

	if err := db.Use(
		dbresolver.Register(dbresolver.Config{}).
			SetMaxOpenConns(maxOpenConns).
			SetMaxIdleConns(maxIdleConns).
			SetConnMaxIdleTime(time.Duration(maxIdleTime) * time.Second),
	); err != nil {
		logger.Panic(err)
	}

	if cfg.ShowStats {
		go func() {
			m, _ := db.DB()
			ticker := time.NewTicker(5 * time.Second)

			for {
				<-ticker.C
				logger.Debugf("current %s conn stats: %#v", cfg.DBType, m.Stats())
			}
		}()
	}

	logger.Infof("connected to %s  [%s - %s]", cfg.DBType, cfg.Addr, cfg.DBName)

	return db
}

func (c *Client) migrateUp(ctx context.Context, cfg Config, sourcePath string) error {
	logger := cmlog.Extract(ctx)

	if cfg.CreateDB {
		db := c.openDB(ctx, cfg, false)
		if err := db.Exec(fmt.Sprintf("CREATE DATABASE %s;", QuoteTo(db, cfg.DBName))).Error; err != nil {
			switch e := err.(type) {
			case *pgconn.PgError:
				// duplicate_database
				if e.Code != "42P04" {
					return err
				}
			case *mysql.MySQLError:
				// ER_DB_CREATE_EXISTS
				if e.Number != 1007 {
					return err
				}
			default:
				return err
			}
		}
	}

	dsn, err := cfg.GetDSN(true)
	if err != nil {
		return err
	}

	stripScheme := func(dsn string, dbType DBType) string {
		return strings.TrimPrefix(dsn, fmt.Sprintf("%s://", dbType))
	}

	m, err := migrate.New(fmt.Sprintf("file://%v", sourcePath), fmt.Sprintf("%s://%s", cfg.DBType, stripScheme(dsn, cfg.DBType)))
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}

	logger.Infof("%s migrations succeed", cfg.DBType)

	return nil
}

func (c *Client) getDSN(ctx context.Context, cfg ConfigGetter, useDB bool) string {
	logger := cmlog.Extract(ctx)

	dsn, err := cfg.GetDSN(useDB)
	if err != nil {
		logger.Panicf("get dsn failed: %v", err)
	}

	return dsn
}

func (c *Client) getDBOpener(ctx context.Context, cfg ConfigGetter) func(string) gorm.Dialector {
	logger := cmlog.Extract(ctx)

	dialector, err := cfg.GetDialector()
	if err != nil {
		logger.Panicf("get dialector failed: %v", err)
	}

	return dialector
}
