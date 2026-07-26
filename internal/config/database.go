package config

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/facily-tech/go-core/env"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	Prefix = "DB_"
)

type ConfigDB struct {
	Engine      string        `env:"ENGINE"`
	DSN         string        `env:"DSN"`
	MaxIdleConn int           `env:"MAX_IDLE,default=5"`
	MaxOpenConn int           `env:"MAX_OPEN,default=10"`
	MaxIdleTime time.Duration `env:"MAX_IDLE_DURATION,default=1m"`
	MaxLifetime time.Duration `env:"MAX_LIFETIME_DURATION,default=5m"`
}

type slogWriter struct{}

func (slogWriter) Printf(format string, v ...any) {
	slog.Default().Info("gorm", slog.String("msg", format), slog.Any("args", v))
}

func InitDB(ctx context.Context) (*ConfigDB, *gorm.DB, error) {
	var config ConfigDB
	if err := env.LoadEnv(ctx, &config, Prefix); err != nil {
		return nil, nil, errors.Wrap(err, "cannot load db environment variable")
	}

	var err error

	var db *gorm.DB

	switch config.Engine {
	case "sqlite":
		db, err = openSqlite(config)
	case "postgres":
		db, err = openPostgres(config)
	default:
		err = errors.Errorf("invalid database engine: %s", config.Engine)
	}

	return &config, db, errors.Wrapf(err, "%s cannot open database", config.Engine)

}

func openSqlite(config ConfigDB) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(config.DSN), &gorm.Config{})
}

func openPostgres(config ConfigDB) (*gorm.DB, error) {
	datasource := config.DSN

	cfg, err := pgxpool.ParseConfig(datasource)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	constr := stdlib.RegisterConnConfig(cfg.ConnConfig)

	sqlDB, err := sql.Open("pgx", constr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sqlDB.SetConnMaxIdleTime(config.MaxIdleTime)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)

	sqlDB.SetMaxIdleConns(config.MaxIdleConn)
	sqlDB.SetMaxOpenConns(config.MaxOpenConn)

	gormLog := logger.New(
		slogWriter{},
		logger.Config{
			SlowThreshold:             400 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
		},
	)

	options := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		FullSaveAssociations:                     true,
		SkipDefaultTransaction:                   true,
		Logger:                                   gormLog,
	}

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), options)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return db, nil
}
