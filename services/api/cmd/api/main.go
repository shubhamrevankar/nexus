package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nexus/api/internal/httpserver"
	"github.com/nexus/api/internal/identity"
	"github.com/nexus/api/internal/platform/config"
	"github.com/nexus/api/internal/platform/database"
	"github.com/nexus/api/internal/platform/logging"
	"github.com/nexus/api/internal/tenancy"
)

func main() {
	appConfig := config.Load()
	logger := logging.New(appConfig.LogLevel)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.Open(ctx, appConfig.DatabaseURL)
	if err != nil {
		logger.Error("database connection failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	if err := database.Migrate(ctx, db, appConfig.MigrationsDir); err != nil {
		logger.Error("database migration failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	identityRepository := identity.NewRepository(db)
	tenancyRepository := tenancy.NewRepository(db)

	server := &http.Server{
		Addr:              ":" + appConfig.Port,
		Handler:           httpserver.NewRouter(logger, identityRepository, tenancyRepository, appConfig.AllowedOrigins, time.Duration(appConfig.SessionTTLHours)*time.Hour),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("api server starting", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("api server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("api server shutdown failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("api server stopped")
}
