# Chirpy

A small Go web server for the Chirpy project, organized in layers similar to an Express application.

## Project structure

```text
.
├── main.go                      # Starts the server (like server.js)
├── internal/
│   ├── app/
│   │   └── app.go               # Wires dependencies together (like app.js)
│   ├── config/
│   │   └── config.go            # Loads .env / environment variables
│   ├── router/
│   │   └── router.go            # Route table: URL + method -> handler
│   ├── handlers/                # Controllers (one file per feature)
│   │   ├── handler.go           # Shared dependencies for all handlers
│   │   ├── health.go
│   │   ├── users.go
│   │   ├── chirps.go
│   │   └── admin.go
│   ├── auth/
│   │   └── auth.go              # Password hashing and verification
│   ├── middleware/
│   │   └── metrics.go           # Request-counting middleware
│   ├── respond/
│   │   └── respond.go           # JSON and error response helpers
│   └── database/                # sqlc-generated database code
├── sql/
│   ├── queries/                 # SQL used to generate database functions
│   └── schema/                  # Goose migrations
└── assets/
```

## How a request flows

```text
main.go
  -> app.New          (build dependencies)
    -> router.New     (register routes)
      -> handlers.X   (run business logic)
        -> respond    (write JSON / errors)
```

## Express comparison

| Express concept        | Chirpy file                     |
| ---------------------- | ------------------------------- |
| `server.js`            | `main.go`                       |
| `app.js`               | `internal/app/app.go`           |
| `app.use()` / routes   | `internal/router/router.go`     |
| Controllers            | `internal/handlers/*.go`        |
| Middleware             | `internal/middleware/*.go`      |
| `res.json()` helpers   | `internal/respond/respond.go`   |
| `config` / `.env`      | `internal/config/config.go`     |
| Database models/queries| `internal/database/*.go`        |

## Running

```bash
go run .
```

The server listens on `:8080` by default. Set `ADDRESS` to use a different
address:

```bash
ADDRESS=:3000 go run .
```

Required environment variables:

- `DB_URL` - PostgreSQL connection string
- `PLATFORM` - set to `dev` to enable `POST /admin/reset`

## Endpoints

| Method | Path             | Description                          |
| ------ | ---------------- | ------------------------------------ |
| GET    | `/api/healthz`   | Health check                         |
| POST   | `/api/users`     | Create a user                        |
| POST   | `/api/login`     | Log in with an email and password    |
| POST   | `/api/refresh`   | Exchange a refresh token for a JWT   |
| POST   | `/api/revoke`    | Revoke a refresh token               |
| POST   | `/api/chirps`    | Create a chirp                       |
| GET    | `/admin/metrics` | View request count                   |
| POST   | `/admin/reset`   | Reset metrics and delete users (dev) |
| GET    | `/app/`          | Serve static files                   |
| GET    | `/assets/`       | Serve assets                         |

## Development checks

```bash
gofmt -w .
go build ./...
go vet ./...
```
