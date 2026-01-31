# URL Shortener with HTMX + DaisyUI

A full-stack URL shortener built with [Encore.go](https://encore.dev), [HTMX](https://htmx.org), and [DaisyUI](https://daisyui.com). Demonstrates interactive server-rendered UIs in Go without a JavaScript framework.

## Architecture

The app has two services:

- **`url`** — Core URL shortening logic with a PostgreSQL database.
- **`frontend`** — Serves the web UI using Go templates, HTMX for interactivity, and DaisyUI for styling.

```
go-htmx-daisyui/
├── encore.app
├── url/
│   ├── url.go              # API endpoints (shorten, list, get)
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

## API Endpoints

| Method | Path            | Description                         |
| ------ | --------------- | ----------------------------------- |
| POST   | `/api/url`      | Shorten a URL                       |
| GET    | `/api/url`      | List all shortened URLs             |
| GET    | `/api/url/:id`  | Get original URL by short ID        |
| GET    | `/`             | Dashboard                           |
| GET    | `/urls`         | URL management page                 |
| POST   | `/htmx/shorten` | Shorten URL (returns HTML fragment) |
| GET    | `/htmx/urls`    | URL list (returns HTML fragment)    |
| GET    | `/static/*path` | Static assets (HTMX, DaisyUI)      |

## Prerequisites

- [Encore CLI](https://encore.dev/docs/install)
- Go 1.24+

## Running Locally

```bash
encore run
```

The app will be available at http://localhost:4000.

## Running Tests

```bash
encore test ./...
```

## How It Works

The frontend uses HTMX to send requests to the backend and swap HTML fragments into the page without full reloads:

1. The URL form uses `hx-post="/htmx/shorten"` to submit and insert a new table row.
2. The URL list uses `hx-get="/htmx/urls"` with `hx-trigger="load"` to populate on page load.
3. The backend renders Go templates and returns HTML fragments directly — no JSON serialization or client-side rendering needed.

## Deployment

Deploy to Encore's cloud with:

```bash
git push encore
```
