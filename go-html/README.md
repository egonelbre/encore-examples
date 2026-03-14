# Message Board — Encore + html/template

A server-rendered message board built with Encore.go and Go's `html/template`.

Demonstrates:
- Raw API endpoints (`//encore:api public raw`) serving HTML
- `html/template` with a base layout and partials
- Embedded static files (`embed.FS`)
- Session-based auth using Encore's cache (`encore.dev/storage/cache`)
- Access control with role-based levels (public, editor, admin)
- Flash messages via cookies
- Form decoding with `gorilla/schema`

## Prerequisites

**Install Encore:**
- **macOS:** `brew install encoredev/tap/encore`
- **Linux:** `curl -L https://encore.dev/install.sh | bash`
- **Windows:** `iwr https://encore.dev/install.ps1 | iex`

## Run locally

```bash
encore run
```

Then open [http://localhost:4000/](http://localhost:4000/).

## Usage

| Page     | Path     | Access  |
|----------|----------|---------|
| Home     | `/`      | Public  |
| Edit     | `/edit`  | Editor+ |
| Admin    | `/admin` | Admin   |
| Login    | `/login` | Public  |

Test accounts:
- `editor` / `editor` — can edit the message
- `admin` / `admin` — can edit and view admin page

## Testing

```bash
encore test ./...
```
