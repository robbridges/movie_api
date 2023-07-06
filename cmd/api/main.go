package main

import (
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"movie_api/internal/data"
	"movie_api/internal/jsonlog"
	"os"
)

const (
	version = "1.0.0"
)

type config struct {
	port int
	env  string
	db   struct {
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	viper.SetConfigFile("local.env")
	viper.AddConfigPath("./")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("init: %w", err))
	}

	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (dev|stag|prod)")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.Parse()

	cfg.limiter.rps = viper.GetFloat64("LIMITER_RPS")
	cfg.limiter.burst = viper.GetInt("LIMITER_BURST")
	cfg.limiter.enabled = viper.GetBool("LIMITER_ENABLED")

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	app := &application{
		config: cfg,
		logger: logger,
	}

	err := app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
