package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/spf13/viper"
	"time"
)

type PostgressConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMODE  string
}

// Open will open a sql connection with the provided Postgres. Callers will need to ensure it's closed
func (app *application) Open(cfg PostgressConfig) (*sql.DB, error) {
	db, err := sql.Open(
		"pgx",
		cfg.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("error Opening DB: %w", err)
	}

	db.SetMaxOpenConns(app.config.db.maxOpenConns)
	// Set the maximum number of idle connections in the pool. Again, passing a value // less than or equal to 0 will mean there is no limit. db.SetMaxIdleConns(cfg.db.maxIdleConns)

	db.SetMaxIdleConns(app.config.db.maxIdleConns)

	duration, err := time.ParseDuration(app.config.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	// Set the maximum idle timeout.
	db.SetConnMaxIdleTime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Failed to ping the database:", err)
	}

	return db, nil
}

func (cfg PostgressConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMODE)
}

func DefaultPostgesTestConfig() PostgressConfig {
	return PostgressConfig{
		Host:     viper.GetString("TEST_DATABASE_HOST"),
		Port:     viper.GetString("TEST_DATABASE_PORT"),
		User:     viper.GetString("TEST_DATABASE_USER"),
		Password: viper.GetString("TEST_DATABASE_PASSWORD"),
		Database: viper.GetString("TEST_DATABASE"),
		SSLMODE:  "disable",
	}
}
