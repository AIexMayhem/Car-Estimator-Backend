package database

import (
    "fmt"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator struct {
    migrationTool *migrate.Migrate
}

func (m *Migrator) Init(conn *Connection, conf *Config, path string) error {
    driver, err := postgres.WithInstance(conn.DB(), &postgres.Config{})
    if err != nil {
        return fmt.Errorf("cannot create postgres driver: %w", err)
    }
    mt, err := migrate.NewWithDatabaseInstance(
        "file://"+path,
        conf.Driver,
        driver,
    )
    if err != nil {
        return fmt.Errorf("cannot init migration tool: %w", err)
    }
    m.migrationTool = mt
    return nil
}

func (m *Migrator) Apply() error {
    err := m.migrationTool.Up()
    if err == migrate.ErrNoChange {
        fmt.Println("Nothing to apply")
        return nil
    }
    if err != nil {
        return fmt.Errorf("migration up failed: %w", err)
    }
    fmt.Println("Migration(s) applied successfully")
    return nil
}

func (m *Migrator) RollBack(steps int) error {
    err := m.migrationTool.Steps(-steps)
    if err == migrate.ErrNoChange {
        fmt.Println("Nothing to rollback")
        return nil
    }
    if err != nil {
        return fmt.Errorf("rollback failed: %w", err)
    }
    fmt.Println("Rollback complete")
    return nil
}