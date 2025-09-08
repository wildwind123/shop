package migration

import (
	"database/sql"
	"fmt"

	"github.com/go-faster/errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigratorParams holds the parameters required for migration.
type MigratorParams struct {
	Source        string // database connection string
	MigrationPath string // path to migration files
}

// GetMysqlMigrator creates and returns a new migrator for MySQL.
func GetMysqlMigrator(params MigratorParams) (*migrate.Migrate, error) {
	// Open MySQL database connection
	db, err := sql.Open("mysql", params.Source)
	if err != nil {
		return nil, errors.Wrap(err, "can't open mysql connection")
	}

	// Create a new MySQL driver instance
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "can't create mysql driver")
	}

	// Construct and return the migrator
	return migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", params.MigrationPath),
		"mysql",
		driver,
	)
}

// UpMigrate performs database migrations upwards, ignoring ErrNoChange.
func UpMigrate(m *migrate.Migrate) error {
	err := m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(err, "can't migrate")
	}
	return nil
}
