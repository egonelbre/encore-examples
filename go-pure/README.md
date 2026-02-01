# URL Shortener — Pure Go

A full-stack URL shortener using pure Go with [HTMX](https://htmx.org) and [DaisyUI](https://daisyui.com). Demonstrates how to build interactive server-rendered UIs using only the Go standard library and `pgx` — no framework required.

This is a conversion of the [go-htmx-daisyui](../go-htmx-daisyui) example, replacing Encore's API annotations and `sqldb` package with standard `net/http` routing and direct PostgreSQL access via `pgx/v5`.

## Architecture

The app has two packages:

- **`url`** — Core URL shortening logic with direct PostgreSQL access via `pgx/v5`.
- **`frontend`** — Serves the web UI using `net/http.ServeMux`, Go templates, HTMX for interactivity, and DaisyUI for styling.

```
go-pure/
├── main.go                  # Entry point, wires up services
├── url/
│   ├── url.go               # URL service (shorten, list, get)
│   ├── url_test.go          # Tests
│   └── migrations/
│       └── 1_create_tables.up.sql
└── frontend/
    ├── front.go             # HTTP handlers & page data
    ├── templates.go         # Template parsing helpers
    ├── static.go            # Static file serving
    ├── templates/
    │   ├── base.html        # Base layout
    │   ├── dashboard.html   # Dashboard page
    │   ├── urls.html        # URL management page
    │   └── partials/        # Shared template partials
    └── static/
        ├── htmx/
        └── daisyui/
```

## HTTP Routes

| Method | Path             | Description                         |
| ------ | ---------------- | ----------------------------------- |
| GET    | `/`              | Dashboard                           |
| GET    | `/urls`          | URL management page                 |
| POST   | `/htmx/shorten`  | Shorten URL (returns HTML fragment) |
| GET    | `/htmx/urls`     | URL list (returns HTML fragment)    |
| GET    | `/static/*`      | Static assets (HTMX, DaisyUI)      |

## Prerequisites

- Go 1.22+
- PostgreSQL

## Running Locally

Create the database and apply migrations, then:

```bash
go run .
```

The app will be available at http://localhost:8080.

## Running Tests

```bash
go test ./...
```

## How It Works

The frontend uses HTMX to send requests to the backend and swap HTML fragments into the page without full reloads:

1. The URL form uses `hx-post="/htmx/shorten"` to submit and insert a new table row.
2. The URL list uses `hx-get="/htmx/urls"` with `hx-trigger="load"` to populate on page load.
3. The backend renders Go templates and returns HTML fragments directly — no JSON serialization or client-side rendering needed.

All templates and static assets are embedded into the binary at compile time using `//go:embed`.
