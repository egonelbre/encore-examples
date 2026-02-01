package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/egonelbre/encore-example/go-pure/frontend"
	"github.com/egonelbre/encore-example/go-pure/url"
)

func main() {
	err := run(context.Background())
	if err != nil {
		slog.Error("run failed", err)
	}
}

func run(ctx context.Context) error {
	urlService, err := url.Connect(ctx, "postgres://user:password@localhost:5432/url")
	if err != nil {
		return fmt.Errorf("failed to open url service: %w", err)
	}
	defer urlService.Close()

	server := frontend.NewServer(urlService)
	return http.ListenAndServe("127.0.0.1:8080", server)
}
