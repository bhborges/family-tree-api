// Package db has built-in configuration for working
// with some databases used in our ecosystem.
package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4/database"
	postgresmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// Necessary for working with postgres driver.
	_ "github.com/lib/pq"

	// This is required to new relic agent to work with Postegres.
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
)

// PostgresConfig holds all necessary configuration to
// run a PostgreSQL instance.
type PostgresConfig struct {
	Address struct {
		Host     string `split_words:"true" required:"false" default:"localhost"`
		Port     int64  `split_words:"true" required:"false" default:"5432"`
		User     string `split_words:"true" required:"false" default:"postgres"`
		Password string `split_words:"true" required:"false" default:"postgres"`
		Database string `split_words:"true" required:"true"`
	}
}

// ProvidePostgresConfig process the configuration needed
// to run a PostgreSQL database.
func ProvidePostgresConfig(l *zap.Logger) (*PostgresConfig, error) {
	var config PostgresConfig
	if err := envconfig.Process("postgres", &config); err != nil {
		l.Error(ErrPostgresEnvConfig.Error(), zap.Error(err))

		return nil, ErrPostgresEnvConfig
	}

	return &config, nil
}

func dsn(config *PostgresConfig) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Address.Host,
		config.Address.Port,
		config.Address.User,
		config.Address.Password,
		config.Address.Database,
	)
}

// ProvidePostgresDatabase initializes a connection with a PostgreSQL database
// with all necessary configuration.
// It executes migration scripts if a path to files are provided.
//
// NOTE: we are using New Relic drive to support spam by default.
// More details:
// https://pkg.go.dev/github.com/newrelic/go-agent/v3/integrations/nrpq
func ProvidePostgresDatabase(config *PostgresConfig, logger *zap.Logger) (*sql.DB, error) {
	logger.Info("opening connection with PostgreSQL database...", zap.String("host", config.Address.Host))

	sqlDB, err := sql.Open("nrpostgres", dsn(config))
	if err != nil {
		return nil, ErrSQLOpenConn
	}

	return sqlDB, nil
}

// ProvidePostgresGORMDatabase creates a new GORM instance based on an
// existing database connection.
func ProvidePostgresGORMDatabase(config *PostgresConfig, logger *zap.Logger, sqlDB *sql.DB) (*gorm.DB, error) {
	logger.Info("opening GORM connection...")

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

// RegisterPostgresDatabaseOnStopHook is used to register a hook to be executed
// whenever the application stops. This will close *sql.DB connection.
func RegisterPostgresDatabaseOnStopHook(lc fx.Lifecycle, db *sql.DB, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("closing PostgreSQL database connection...")

			err := db.Close()
			if err != nil {
				logger.Error(ErrSQLCloseConn.Error(), zap.Error(err))

				return ErrSQLCloseConn
			}

			return nil
		},
	})
}

// ProvidePostgresMigrateDriver returns a new driver for
// working with migrations in a PostegreSQL database.
func ProvidePostgresMigrateDriver(db *sql.DB, logger *zap.Logger, config *MigrateConfig) (database.Driver, error) {
	if db == nil || config.Path == "" {
		return nil, nil
	}

	driver, err := postgresmigrate.WithInstance(db, &postgresmigrate.Config{})
	if err != nil {
		logger.Error(ErrPostgresMigrateDriver.Error(), zap.Error(err))

		return nil, ErrPostgresMigrateDriver
	}

	return driver, nil
}

// PostgresModule wrapper for uber/fx.
//
//nolint:gochecknoglobals
var PostgresModule = fx.Options(
	fx.Provide(ProvidePostgresConfig),
	fx.Provide(ProvidePostgresMigrateDriver),
	fx.Provide(ProvidePostgresDatabase),
	fx.Provide(ProvidePostgresGORMDatabase),
	fx.Invoke(RegisterPostgresDatabaseOnStopHook),
)
