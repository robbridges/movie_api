package main

import (
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	version = "1.0.0"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	logger *log.Logger
}

func init() {
	viper.SetConfigFile("local.env")
	viper.AddConfigPath("./")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("init: %w", err))
	}
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (dev|stag|prod)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	dbConfig := DefaultPostgesTestConfig()
	db, err := Open(dbConfig)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Printf("Connected to db")

	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)

	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatal(err)
	}
}
