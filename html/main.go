package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/egonelbre/encore-example/go-pure/frontend"
	"github.com/egonelbre/encore-example/go-pure/pgdb"
	"github.com/egonelbre/encore-example/go-pure/url"
)

func main() {
	err := run(context.Background())
	if err != nil {
		slog.Error("run failed", slog.Any("error", err))
	}
}

func run(ctx context.Context) error {
	db, err := pgdb.Connect(ctx, "postgres://user:password@localhost:5432/url")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	urlService := url.NewService(db.URLs())
	server := frontend.NewServer(urlService)

	return http.ListenAndServe("127.0.0.1:8080", server)
}
