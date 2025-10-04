package postgres_test

import (
	"context"
	"database/sql"
	"embed"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	postgresdriver "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

//go:embed../../../db/migrations/*.sql
var migrationsFS embed.FS

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	pgContianer, err := postgres.Run(
		ctx,
		"postgres:15",
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(10*time.Second),
		),
		testcontainers.WithEnv(map[string]string{
			"POSTGRES_DB":       "auth_test",
			"POSTGRES_USER":     "auth_user",
			"POSTGRES_PASSWORD": "auth_password",
		}),
	)
	require.NoError(t, err)

	connString, err := pgContianer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := sql.Open("postgres", connString)
	require.NoError(t, err)

	runMigrations(t, db)

	cleanup := func() {
		db.Close()
		pgContianer.Terminate(ctx) // Clean up container
	}

	return db, cleanup

}

func runMigrations(t *testing.T, db *sql.DB) {

	driver, err := postgresdriver.WithInstance(db, &postgresdriver.Config{})
	require.NoError(t, err)

	source, err := iofs.New(migrationsFS, "db/migrations")
	require.NoError(t, err)

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	require.NoError(t, err)

	err = m.Up()
	require.NoError(t, err)

	version, dirty, err := m.Version()
	require.NoError(t, err)
	require.False(t, dirty, "migrations should not be dirty")
	t.Logf("Migrations applied up to version %d", version)

}
