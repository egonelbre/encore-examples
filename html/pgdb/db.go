package pgdb

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type Root struct {
	db *pgxpool.Pool
}

// Connect creates a new Root instance and runs pending migrations.
func Connect(ctx context.Context, dsn string) (*Root, error) {
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := runMigrations(ctx, db); err != nil {
		db.Close()
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	return &Root{db: db}, nil
}

// runMigrations applies pending database migrations using tern.
func runMigrations(ctx context.Context, db *pgxpool.Pool) error {
	conn, err := db.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquiring connection: %w", err)
	}
	defer conn.Release()

	migrator, err := migrate.NewMigrator(ctx, conn.Conn(), "public.schema_version")
	if err != nil {
		return fmt.Errorf("creating migrator: %w", err)
	}

	migrations, err := fs.Sub(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("reading migrations: %w", err)
	}

	if err := migrator.LoadMigrations(migrations); err != nil {
		return fmt.Errorf("loading migrations: %w", err)
	}

	return migrator.Migrate(ctx)
}

func (s *Root) Close() error {
	s.db.Close()
	return nil
}
