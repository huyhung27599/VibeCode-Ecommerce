# Ecommerce Backend (Go + Gin)

## Layout
```
backend/
├── cmd/api/              # entrypoint
├── internal/
│   ├── config/           # env-based config
│   ├── server/           # gin engine + route wiring
│   ├── handler/          # HTTP handlers
│   ├── service/          # business logic
│   ├── repository/       # data access (gorm/redis)
│   ├── domain/           # core entities
│   └── middleware/       # request_id, logger, recovery, auth
├── pkg/
│   ├── database/         # postgres + redis clients
│   ├── logger/           # slog setup
│   └── response/         # JSON response envelope
├── migrations/
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

## Quick start
```bash
cp .env.example .env
make docker-up            # start postgres + redis
make tidy
make run
```

Health: `GET /health` · Readiness: `GET /ready` · Ping: `GET /api/v1/ping`

## Adding a resource
1. Define entity in `internal/domain/`
2. Add repository interface + gorm impl in `internal/repository/`
3. Add business logic in `internal/service/`
4. Add handler in `internal/handler/`
5. Wire routes in `internal/server/server.go`
