# Go Boilerplate Checklist

## Done

- [x] Config loading (`.env` → typed config struct with validation)
- [x] Structured logging (slog JSON handler, env-based log level)
- [x] Request logging middleware (method, path, status, duration, IP)
- [x] Panic recovery middleware
- [x] DB connection with ping check
- [x] Migrations (golang-migrate, auto-run on startup)
- [x] Domain-driven folder structure (`internal/<domain>/`)
- [x] Repository → Service → Handler layered architecture
- [x] Container-based DI wiring
- [x] Graceful shutdown (SIGTERM/SIGINT, 5s drain timeout)
- [x] Linting (`.golangci.yml`)
- [x] Makefile (run, build, lint, fmt, tidy, migrate-up, migrate-down)
- [x] `.env.example`

---

## In Progress / Left

### High Priority

- [ ] Fix hardcoded port `8090` — read from config
- [ ] Pass `context.Context` through repo queries (request cancellation, timeouts)
- [ ] DB connection pool config (`SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`)
- [ ] Interfaces for service and repository (decoupling, testability)
- [ ] Standardized error response type (`{ error, message, statusCode }`)
- [ ] Request body validation (`go-playground/validator`)

### Medium Priority

- [ ] Request ID middleware (UUID per request, attached to logger)
- [ ] Structured `AppError` type (distinguish domain errors from infra errors)
- [ ] Health check returning JSON with DB ping check
- [ ] CORS middleware
- [ ] API versioning prefix (`/api/v1/`)
- [ ] Response envelope helper (`{ data, meta }`)

### Lower Priority

- [ ] Rate limiting middleware
- [ ] Use `shuttingDown` flag in middleware to reject requests during shutdown
- [ ] Close DB connection in graceful shutdown
- [ ] Unit tests
- [ ] Integration tests
- [ ] OpenAPI / Swagger docs
