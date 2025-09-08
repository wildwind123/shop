package main

import (
	"context"
	"flag"
	"github/wildwind123/shop/internal/adminapi"
	"github/wildwind123/shop/internal/provider"
	"github/wildwind123/shop/pkg/config"
	"github/wildwind123/shop/pkg/db"
	"github/wildwind123/shop/pkg/migration"
	"github/wildwind123/shop/pkg/ogenapi"
	"log/slog"
	"net/http"
	"os"

	"github.com/wildwind123/slogger"
	"github.com/wildwind123/xutils"
)

func main() {

	// logger
	logger := slogger.NewLogger(&slogger.Options{
		Level:     slog.LevelDebug,
		AddSource: true,
		Writer:    os.Stdout,
		App:       "shop",
		Build:     "v1.2",
	})
	ctx := context.Background()
	ctx = slogger.ToCtx(ctx, logger)

	// load conig
	configPath := flag.String("config", "../../docs/config/config.yaml", "path of config file")
	flag.Parse()

	logger.Info("config", slog.String("path", *configPath))
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Error("cant load config", slog.Any("err", err))
		return
	}

	// migrator
	m, err := migration.GetMysqlMigrator(migration.MigratorParams{
		Source:        cfg.Database.Source,
		MigrationPath: cfg.Database.MigrationPath,
	})
	if err != nil {
		logger.Error("cant get mysql migrator", slog.Any("err", err))
		return
	}
	err = migration.UpMigrate(m)
	if err != nil {
		logger.Error("cant migrate", slog.Any("err", err))
		return
	}

	// db
	db, err := db.New(cfg.Database.Source)
	if err != nil {
		logger.Error("cant get db", slog.Any("err", err))
		return
	}

	// provider
	pr := provider.Provider{
		DB: db,
	}

	// run server
	// run admin api
	adminApi := adminapi.Adminapi{
		ProductHandler: &adminapi.ProductHandler{
			Provider: &pr,
		},
	}
	muxServer := http.NewServeMux()

	oa, err := ogenapi.NewServer(adminApi, ogenapi.WithPathPrefix("/api"))
	if err != nil {
		logger.Error("cant create ogenapi server")
		return
	}

	muxServer.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		ctx := xutils.RequestToCtx(ctx, r)
		r = r.WithContext(ctx)
		oa.ServeHTTP(w, r)
	})

	logger.Info("run admin api", slog.String("addr", cfg.Server.Addr))
	err = http.ListenAndServe(cfg.Server.Addr, muxServer)
	if err != nil {
		logger.Error("cant run server", slog.Any("err", err))
	}
}
