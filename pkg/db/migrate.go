package db

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"
	"go.uber.org/zap"

	// This is required for migration to work.
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateConfig holds all necessary configuration
// to run a database migration.
type MigrateConfig struct {
	Path        string `split_words:"true" required:"false"`
	EnableDebug bool   `split_words:"true" required:"false" default:"false"`
}

// ProvideMigrateConfig process the configuration needed
// to run a database migration.
func ProvideMigrateConfig(l *zap.Logger) (*MigrateConfig, error) {
	var config MigrateConfig
	if err := envconfig.Process("migrate", &config); err != nil {
		l.Error(ErrMigrateEnvConfig.Error(), zap.Error(err))

		return nil, ErrMigrateEnvConfig
	}

	return &config, nil
}

// ProvideMigrateParams holds all parametes from ProvideMigrate func.
type ProvideMigrateParams struct {
	fx.In

	Config *MigrateConfig
	Logger *zap.Logger
	Driver database.Driver `optional:"true"`
}

// ProvideMigrate configures a new migrate instance from the provided driver.
func ProvideMigrate(p ProvideMigrateParams) (*migrate.Migrate, error) {
	if p.Driver == nil {
		return nil, nil
	}

	m, err := migrate.NewWithDatabaseInstance(p.Config.Path, "postgres", p.Driver)
	if err != nil {
		p.Logger.Error(ErrMigrateNewInstance.Error(), zap.Error(err))

		return nil, ErrMigrateNewInstance
	}

	m.Log = &migrateZapLogger{logger: p.Logger, config: p.Config}

	return m, nil
}

// ExecuteMigrationParams holds all parametes from ExecuteMigration func.
type ExecuteMigrationParams struct {
	fx.In

	M      *migrate.Migrate `optional:"true"`
	Logger *zap.Logger
}

// ExecuteMigration runs all script migrations UP.
func ExecuteMigration(e ExecuteMigrationParams) error {
	if e.M == nil {
		e.Logger.Warn("no migration path provided, skipping...")

		return nil
	}

	migrationVersion, isDirty, err := e.M.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		e.Logger.Error(ErrMigrateVersion.Error(), zap.Error(err))

		return ErrMigrateVersion
	}

	e.Logger.Info("running migration scripts...",
		zap.Uint("migration_version", migrationVersion),
		zap.Bool("is_dirty", isDirty),
	)

	err = e.M.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		e.Logger.Info("no migrations to be executed, already up to date!")

		return nil
	}

	if err != nil {
		e.Logger.Error(ErrMigrateApply.Error(), zap.Error(err))

		return ErrMigrateApply
	}

	return nil
}

// migrateZapLogger implements migrate.Logger interface.
type migrateZapLogger struct {
	logger *zap.Logger
	config *MigrateConfig
}

// Printf implementation from migrate.Logger interface.
func (m *migrateZapLogger) Printf(format string, v ...interface{}) {
	m.logger.Info(fmt.Sprintf(format, v...))
}

// Verbose implementation from migrate.Logger interface.
func (m *migrateZapLogger) Verbose() bool {
	return m.config.EnableDebug
}

// MigrateModule wrapper for uber/fx.
//
//nolint:gochecknoglobals
var MigrateModule = fx.Options(
	fx.Provide(ProvideMigrateConfig),
	fx.Provide(ProvideMigrate),
	fx.Invoke(ExecuteMigration),
)
