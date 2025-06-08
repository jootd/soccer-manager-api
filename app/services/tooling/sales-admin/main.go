package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jootd/soccer-manager/app/services/tooling/sales-admin/commands"
	"github.com/jootd/soccer-manager/business/sdk/sqldb"
	"go.uber.org/zap"
)

var build = "develop"

type Config struct {
	User         string
	Password     string
	Host         string
	Name         string
	MaxIdleConns int
	MaxOpenConns int
	DisableTLS   bool
}

func main() {
	if err := run(zap.NewNop().Sugar()); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	// Load config from environment manually (since we're not using conf.Parse)
	cfg := Config{
		User:         getEnv("DB_USER", "myuser"),
		Password:     getEnv("DB_PASSWORD", "mypass"),
		Host:         getEnv("DB_HOST", "localhost"),
		Name:         getEnv("DB_NAME", "soccer"),
		MaxIdleConns: getEnvInt("DB_MAXIDLECONNS", 2),
		MaxOpenConns: getEnvInt("DB_MAXOPENCONNS", 10),
		DisableTLS:   getEnvBool("DB_DISABLETLS", true),
	}

	if len(os.Args) < 2 {
		printUsage()
		return commands.ErrHelp
	}

	return processCommands(os.Args[1], log, cfg)
}

func processCommands(cmd string, log *zap.SugaredLogger, cfg Config) error {
	dbConfig := sqldb.Config{
		User:         cfg.User,
		Password:     cfg.Password,
		Host:         cfg.Host,
		Name:         cfg.Name,
		MaxIdleConns: cfg.MaxIdleConns,
		MaxOpenConns: cfg.MaxOpenConns,
		DisableTLS:   cfg.DisableTLS,
	}

	switch cmd {
	case "migrate":
		if err := commands.Migrate(dbConfig); err != nil {
			return fmt.Errorf("migrating database: %w", err)
		}
	case "seed":
		if err := commands.Seed(dbConfig); err != nil {
			return fmt.Errorf("seeding database: %w", err)
		}
	default:
		printUsage()
		return commands.ErrHelp
	}

	return nil
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  migrate    Create the schema in the database")
	fmt.Println("  seed       Add data to the database")
}

// --- Helper functions to read env vars with fallback and parsing ---
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return fallback
}
