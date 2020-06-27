package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Fs02/go-todo-backend/api"
	"github.com/Fs02/rel"
	"github.com/Fs02/rel/adapter/mysql"
	_ "github.com/go-sql-driver/mysql" // TODO: use postgres
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewProduction(zap.Fields(zap.String("type", "main")))
	shutdowns []func() error
)

func main() {
	var (
		ctx        = context.Background()
		port       = os.Getenv("PORT")
		repository = initRepository()
		mux        = api.NewMux(repository)
		server     = http.Server{
			Addr:    ":" + port,
			Handler: mux,
		}
		shutdown = make(chan struct{})
	)

	go gracefulShutdown(ctx, &server, shutdown)

	logger.Info("server starting: http://localhost" + server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal("server error", zap.Error(err))
	}

	<-shutdown
}

func initRepository() rel.Repository {
	var (
		logger, _ = zap.NewProduction(zap.Fields(zap.String("type", "repository")))
		dsn       = fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			os.Getenv("MYSQL_USERNAME"),
			os.Getenv("MYSQL_PASSWORD"),
			os.Getenv("MYSQL_HOST"),
			os.Getenv("MYSQL_PORT"),
			os.Getenv("MYSQL_DATABASE"))
	)

	adapter, err := mysql.Open(dsn)
	if err != nil {
		logger.Fatal(err.Error(), zap.Error(err))
	}
	// add to graceful shutdown list.
	shutdowns = append(shutdowns, adapter.Close)

	repository := rel.New(adapter)
	repository.Instrumentation(func(ctx context.Context, op string, message string) func(err error) {
		// no op for rel functions.
		if strings.HasPrefix(op, "rel-") {
			return func(error) {}
		}

		t := time.Now()

		return func(err error) {
			duration := time.Since(t)
			if err != nil {
				logger.Error(message, zap.Error(err), zap.Duration("duration", duration), zap.String("operation", op))
			} else {
				logger.Info(message, zap.Duration("duration", duration), zap.String("operation", op))
			}
		}
	})

	return repository
}

func gracefulShutdown(ctx context.Context, server *http.Server, shutdown chan struct{}) {
	var (
		sigint = make(chan os.Signal, 1)
	)

	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	logger.Info("shutting down server gracefully")

	// stop receiving any request.
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("shutdown error", zap.Error(err))
	}

	// close any other modules.
	for i := range shutdowns {
		shutdowns[i]()
	}

	close(shutdown)
}
