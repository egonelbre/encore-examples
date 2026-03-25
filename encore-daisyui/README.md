# URL Shortener with DaisyUI

A full-stack URL shortener built with [Encore.go](https://encore.dev), [DaisyUI](https://daisyui.com), and [Tailwind CSS](https://tailwindcss.com). Demonstrates server-rendered UIs in Go using HTML templates with DaisyUI components.

## Architecture

The app has two services:

- **`url`** — Core URL shortening logic with a PostgreSQL database.
- **`frontend`** — Serves the web UI using Go templates, with DaisyUI and Tailwind CSS for styling.

```
go-daisyui/
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
        └── daisyui/         # DaisyUI CSS & Tailwind browser JS
```

## API Endpoints

| Method | Path             | Description                     |
| ------ | ---------------- | ------------------------------- |
| POST   | `/api/url`       | Shorten a URL                   |
| GET    | `/api/url`       | List all shortened URLs         |
| GET    | `/api/url/:id`   | Get original URL by short ID    |
| GET    | `/`              | Dashboard                       |
| GET    | `/urls`          | URL management page             |
| POST   | `/url`           | Shorten URL (redirects to urls) |
| GET    | `/static/*path`  | Static files (CSS, JS)          |

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

## Deployment

Deploy to Encore's cloud with:

```bash
git push encore
```
