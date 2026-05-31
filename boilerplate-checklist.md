# Go Boilerplate Checklist

## Done

- [x] Config loading (`.env` → typed config struct with validation)
- [x] Structured logging (slog JSON handler, env-based log level)
- [x] Request logging middleware (method, path, status, duration, IP, request ID)
- [x] Panic recovery middleware
- [x] DB connection with ping check
- [x] DB connection pool config (`SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`)
- [x] Migrations (golang-migrate, auto-run on startup)
- [x] Domain-driven folder structure (`internal/<domain>/`)
- [x] Repository → Service → Handler layered architecture
- [x] Interfaces for service and repository (decoupling, testability)
- [x] Container-based DI wiring
- [x] `context.Context` passed through repo queries (request cancellation + timeout middleware)
- [x] Graceful shutdown (SIGTERM/SIGINT, 5s drain timeout, shuttingDown middleware)
- [x] Standardized error response (`{ error, code, status }`)
- [x] `AppError` type (NotFound, Unauthorized, BadRequest, Forbidden, Conflict, Internal)
- [x] Response envelope (`{ data, meta }`)
- [x] Request body validation setup (`go-playground/validator`)
- [x] Request ID middleware (UUID per request, logged + returned in response header)
- [x] Health check with DB ping returning JSON
- [x] CORS middleware
- [x] API versioning prefix (`/api/v1/`)
- [x] Rate limiting middleware
- [x] Linting (`.golangci.yml`)
- [x] Makefile (run, build, lint, fmt, tidy, migrate-up, migrate-down)
- [x] `.env.example`

---

## Left

- [ ] Unit tests
- [ ] Integration tests
- [ ] OpenAPI / Swagger docs
