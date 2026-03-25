package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/egonelbre/encore-example/go-pure/frontend"
	"github.com/egonelbre/encore-example/go-pure/pgdb"
	"github.com/egonelbre/encore-example/go-pure/url"
)

func main() {
	dev := flag.Bool("dev", false, "enable development mode with auto-reload")
	flag.Parse()

	err := run(context.Background(), *dev)
	if err != nil {
		slog.Error("run failed", slog.Any("error", err))
	}
}

func run(ctx context.Context, dev bool) error {
	db, err := pgdb.Connect(ctx, "postgres://user:password@localhost:5432/url")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	var templates *frontend.Templates
	if dev {
		templates = frontend.NewDevTemplates("frontend")
	} else {
		templates = frontend.NewTemplates()
	}

	urlService := url.NewService(db.URLs())
	server := frontend.NewServer(urlService, templates)

	return http.ListenAndServe("127.0.0.1:8080", server)
}
