package main

import (
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"movie_api/internal/data"
	"movie_api/internal/jsonlog"
	"movie_api/internal/mailer"
	"os"
	"strings"
	"sync"
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
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
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

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	flag.Parse()

	cfg.limiter.rps = viper.GetFloat64("LIMITER_RPS")
	cfg.limiter.burst = viper.GetInt("LIMITER_BURST")
	cfg.limiter.enabled = viper.GetBool("LIMITER_ENABLED")

	cfg.smtp.host = viper.GetString("EMAIL_HOST")
	cfg.smtp.port = viper.GetInt("EMAIL_PORT")
	cfg.smtp.username = viper.GetString("EMAIL_USERNAME")
	cfg.smtp.password = viper.GetString("EMAIL_PASSWORD")
	cfg.smtp.sender = "support@moviebuffs.com"

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	app := &application{
		config: cfg,
		logger: logger,
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err := app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
