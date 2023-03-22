package db

import "errors"

var (
	// ErrMigrateEnvConfig is returned if some error occurs setting up the environent vars.
	ErrMigrateEnvConfig = errors.New("migration: unable to setup environment variables")
	// ErrMigrateNewInstance is returned if not possible to instantiate migration.
	ErrMigrateNewInstance = errors.New("error creating new instance of migration")
	// ErrMigrateVersion is returned if an error occurs while checking migration version.
	ErrMigrateVersion = errors.New("error checking migration version")
	// ErrMigrateApply is returned if unable to apply migration scripts.
	ErrMigrateApply = errors.New("unexpected error applying migration scripts")

	// ErrSQLOpenConn is returned if unable to open connection with SQL database.
	ErrSQLOpenConn = errors.New("unable to open connection with SQL database")
	// ErrSQLCloseConn is returned if unable to close connection with SQL database.
	ErrSQLCloseConn = errors.New("unable to close connection with SQL database")

	// ErrPostgresEnvConfig is returned if some error occurs setting up the environent vars.
	ErrPostgresEnvConfig = errors.New("postgres: unable to setup environment variables")
	// ErrPostgresMigrateDriver is returned if unable to instantiate postgres migration driver.
	ErrPostgresMigrateDriver = errors.New("unable to instantiate postgres migration driver")
)
