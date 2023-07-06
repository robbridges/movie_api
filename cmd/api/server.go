package main

import (
	"context"
	"errors"
	"fmt"
	"movie_api/internal/data"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutDownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		// listen for sigint or sigterm calls, relay them to the channel
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		app.logger.PrintInfo("Shutting down server", map[string]string{
			"signal": s.String(),
		})

		// in flight requests have a 30 second grace period
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		shutDownError <- srv.Shutdown(ctx)
	}()

	dbConfig := DefaultPostgesTestConfig()

	db, err := app.Open(dbConfig)
	if err != nil {
		app.logger.PrintFatal(err, nil)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Failed to ping the database:", err)
	}

	app.models = data.NewModels(db)

	app.logger.PrintInfo("Connected to db", nil)

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	err = srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutDownError
	if err != nil {
		return err
	}

	app.logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
