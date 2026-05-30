package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log/slog"

	"github.com/icecarti/feed_service/database"
	"github.com/icecarti/feed_service/grpc_server"
	"github.com/icecarti/feed_service/services"
	"github.com/joho/godotenv"
)

func setupLogger(env, filename string) (*slog.Logger, error) {
    var out io.Writer
    if filename != "" {
        f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
        if err != nil {
            return nil, fmt.Errorf("unable to open log file: %w", err)
        }
        out = f
    } else {
        out = os.Stdout
    }

    lvl := slog.LevelInfo
    if env == "local" {
        lvl = slog.LevelDebug
    }

    handler := slog.NewTextHandler(out, &slog.HandlerOptions{
        AddSource: true,
        Level:     lvl,
    })
    return slog.New(handler), nil
}

func main() {
    if err := godotenv.Load(); err != nil {
        fmt.Fprintf(os.Stderr, "DEBUG: no .env loaded: %v\n", err)
    } else {
        fmt.Println("DEBUG: .env loaded")
    }

    env := os.Getenv("ENV")
    logFile := os.Getenv("LOG_FILE")
    dbAddr := os.Getenv("DB_ADDR")
    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    fmt.Printf("DEBUG: ENV=%q LOG_FILE=%q DB_ADDR=%q DB_USER=%q DB_NAME=%q\n",
        env, logFile, dbAddr, dbUser, dbName)

    logger, err := setupLogger(env, logFile)
    if err != nil {
        fmt.Fprintf(os.Stderr, "FATAL: logger setup failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("DEBUG: logger initialized")

    logger.Info("Logger initialized", "env", env, "log_file", logFile)

    cfg := &database.Config{
        Driver:   "postgres",
        Addr:     dbAddr,
        User:     dbUser,
        Password: dbPass,
        DBName:   dbName,
    }
    var conn database.Connection
    if err := conn.Init(cfg); err != nil {
        fmt.Fprintf(os.Stderr, "FATAL: DB init failed: %v\n", err)
        os.Exit(1)
    }
    defer conn.Close()
    fmt.Println("DEBUG: database connected")

    logger.Info("Database connected", "addr", dbAddr, "db_name", dbName)

    migr := &database.Migrator{}
    if err := migr.Init(&conn, cfg, "database/migrations"); err != nil {
        fmt.Fprintf(os.Stderr, "FATAL: migrator init failed: %v\n", err)
        os.Exit(1)
    }
    if err := migr.Apply(); err != nil {
        fmt.Fprintf(os.Stderr, "FATAL: migrations apply failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("DEBUG: migrations applied")

    logger.Info("Migrations applied successfully")

    listingRepo := database.NewListingRepository(conn.DB())
    favoriteRepo := database.NewFavoriteRepository(conn.DB())
    fmt.Println("DEBUG: repositories created")

    feedSvc := services.NewFeedService(listingRepo, favoriteRepo)
    fmt.Println("DEBUG: feed service created")

    grpcSrv := grpc_server.NewGRPCServer(":50052", feedSvc, logger)
    go func() {
        if err := grpcSrv.Start(); err != nil {
            fmt.Fprintf(os.Stderr, "FATAL: gRPC server failed: %v\n", err)
            os.Exit(1)
        }
    }()
    fmt.Println("DEBUG: gRPC server started on :50052")
    logger.Info("FeedService is running on :50052")

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit

    fmt.Println("DEBUG: shutdown signal received")
    logger.Info("Shutdown signal received; stopping gRPC server...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    grpcSrv.Stop()
    <-ctx.Done()

    fmt.Println("DEBUG: gRPC server stopped. Exiting.")
    logger.Info("gRPC server stopped. Exiting.")
}
