package database

import (
    "database/sql"
    "fmt"

    _ "github.com/lib/pq"
)

type Connection struct {
    db *sql.DB
}

type Config struct {
    Driver   string
    Addr     string
    User     string
    Password string
    DBName   string
}

func (conf *Config) getConnString(defaultConn bool) string {
    dbName := conf.DBName
    if defaultConn {
        dbName = ""
    }
    return fmt.Sprintf(
        "%s://%s:%s@%s/%s?sslmode=disable",
        conf.Driver, conf.User, conf.Password, conf.Addr, dbName,
    )
}

func CreateDBIfNotExists(conf *Config) error {
    defaultConn, err := sql.Open(conf.Driver, conf.getConnString(true))
    if err != nil {
        return fmt.Errorf("invalid default connection args: %w", err)
    }
    defer defaultConn.Close()

    var exists bool
    if err := defaultConn.
        QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", conf.DBName).
        Scan(&exists); err != nil {
        return fmt.Errorf("database existence check failed: %w", err)
    }

    if !exists {
        if _, err := defaultConn.Exec("CREATE DATABASE " + conf.DBName); err != nil {
            return fmt.Errorf("database creation failed: %w", err)
        }
    }
    return nil
}

func (c *Connection) Init(conf *Config) error {
    if err := CreateDBIfNotExists(conf); err != nil {
        return err
    }
    db, err := sql.Open(conf.Driver, conf.getConnString(false))
    if err != nil {
        return fmt.Errorf("invalid connection args: %w", err)
    }
    if err := db.Ping(); err != nil {
        return fmt.Errorf("database ping failed: %w", err)
    }
    c.db = db
    return nil
}

func (c *Connection) Close() {
    if c.db != nil {
        c.db.Close()
    }
}

func (c *Connection) DB() *sql.DB {
    return c.db
}