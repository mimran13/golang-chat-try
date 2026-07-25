# 45-Day Go Production Mastery Plan

**Goal:** In 45 days, go from "solid Go boilerplate" to "I can build, ship, and operate a production-grade Go service and I understand every layer of it."

**North star deliverable:** `go-chat` deployed to a real cloud environment, with WebSockets, multi-instance scaling, observability, CI/CD, secrets, load-tested, and hardened. Plus a personal notes doc per phase.

**Reference docs in this repo:**
- [learning-checklist.html](./learning-checklist.html) — tick items off as you cover them
- [boilerplate-checklist.md](../boilerplate-checklist.md) — foundation, mostly done
- Each numbered doc (`01-*.md` … `16-*.md`) — running knowledge log

---

## Ground rules

1. **Daily rhythm** (~2.5–3 hours). Adjust honestly if a day is shorter.
    - **20 min** — read one concept (from checklist or resources below)
    - **90 min** — code the day's deliverable
    - **20 min** — write a `docs/notes-YYYY-MM-DD.md` (what you learned, gotchas, questions)
    - **10 min** — tick items on `learning-checklist.html`
2. **Never skip the notes.** They are the compound interest. If you don't write, you don't remember.
3. **One commit per topic minimum.** Not one giant commit at end of week. Small commits show the journey.
4. **When stuck > 30 min, ask.** Don't burn a day chasing a shadow. Write down what you tried, then ask.
5. **Weekly review Sundays.** Re-open checklist, re-open notes. What's shaky? Redo those.

---

## Phase 1 — Go language fluency (Days 1–8)

**Why first:** the chat app will beat concurrency into you; but the *language basics* (interfaces, errors, slices, strings, generics) you have to seek out. Do them first while your fingers are fresh.

### Day 1 — Basics, values, types
- Zero values for every type; write a tiny program that prints them all.
- Value vs pointer semantics: write a `Counter` struct, method with value receiver breaks, pointer receiver fixes.
- Type conversion vs type assertion — a program that does both.
- **Deliverable:** `scratch/day01_types.go` + notes.

### Day 2 — Strings, bytes, runes
- Iterate a UTF-8 string with `for i, r := range s`; contrast with `s[i]`.
- Convert `[]byte` ↔ `string` and measure allocations with `-benchmem`.
- Use `strings.Builder`.
- **Deliverable:** `scratch/day02_strings_test.go` with a benchmark.

### Day 3 — Slices, maps, arrays
- Print a slice header (ptr, len, cap) using `unsafe` or just observation via `cap()`.
- Reproduce the **slice aliasing gotcha**: `append` to a sub-slice mutates the parent.
- Reproduce **map random iteration order**.
- Preallocate: benchmark `make([]int, 0, N)` vs `make([]int, 0)` with growing appends.
- **Deliverable:** `scratch/day03_slice_gotchas_test.go`.

### Day 4 — Constants, iota, enums
- Model chat "message types" as a `type MessageType int` with `iota`.
- Add a `String()` method to make it a `Stringer`.
- Add flag-style `iota` (bit-shift) for a `RoomPermission` bitmask.
- **Deliverable:** commit to `internal/message/type.go` (start the real chat domain here).

### Day 5 — Interfaces deep dive
- Consumer-defined interfaces (already applied in the repo — study `internal/user/service.go`).
- The nil-interface-vs-interface-holding-nil trap: write a function that returns `*AppError` as `error`; show the trap.
- Use `io.Reader` / `io.Writer` in a small program (e.g. a hash-computing helper).
- Compile-time interface check: `var _ userRepository = (*Repository)(nil)`.
- **Deliverable:** Add compile-time checks to the actual repo interfaces. Notes on the nil-interface trap.

### Day 6 — Errors mastered
- Add error wrapping to your repo: change `return err` → `return fmt.Errorf("finding user by username %q: %w", username, err)` in real places.
- Add a `Unwrap()` method to `AppError`.
- Try `errors.Join` (multi-error). Write a validator that returns multiple field errors joined.
- Read: [Go 1.13 errors blog](https://go.dev/blog/go1.13-errors).
- **Deliverable:** the repo gains wrapped errors, tests still pass.

### Day 7 — Control flow, functions, closures
- Type switch: write a function that takes `any` and prints kind + value.
- Named returns + `defer`: write a wrapping function that annotates the returned error on panic.
- First-class functions: refactor one middleware to use a helper that composes middlewares.
- **Deliverable:** `internal/middlewares/compose.go` (utility).

### Day 8 — Generics
- Write a generic `SliceMap[T, U any](in []T, f func(T) U) []U`.
- Write a generic `KeysOf[K comparable, V any](m map[K]V) []K`.
- Read the [Go generics tutorial](https://go.dev/doc/tutorial/generics).
- **Deliverable:** `internal/genericutil/` package used somewhere in the repo.

**End of Phase 1 review (Day 8 evening):** re-open learning-checklist.html. Tick basics/strings/slices/constants/interfaces/errors/generics. If a card still feels shaky, add it to a "revisit" list.

---

## Phase 2 — Concurrency deep dive (Days 9–14)

**Why standalone:** Learn concurrency in tiny isolated programs *before* wiring it into WebSockets. That way when the hub explodes, you know why.

### Day 9 — Goroutines + closure capture gotcha
- Reproduce the classic loop-variable capture bug (using Go < 1.22 semantics via explicit example).
- Fix with variable shadowing / passing as parameter.
- **Deliverable:** `scratch/day09_goroutines_test.go`.

### Day 10 — Channels (unbuffered vs buffered)
- Producer/consumer with unbuffered channel.
- Same with buffered — measure timing difference under contention.
- `range` over channel + `close`. Understand "receive from closed channel" semantics.
- **Deliverable:** notes on when buffering helps vs. hurts.

### Day 11 — Select + nil channel + default
- `select` with 3 cases.
- Set a case's channel to `nil` to disable it dynamically (a real pattern).
- `default:` for non-blocking send/receive.
- **Deliverable:** a small "event router" program.

### Day 12 — Sync primitives
- `sync.Mutex` protecting a counter — reproduce a race, then fix it. Run with `-race`.
- `sync.RWMutex` — writer vs many readers benchmark.
- `sync.WaitGroup` for a fan-out of N workers.
- `sync.Once` for lazy init.
- **Deliverable:** benchmarks committed.

### Day 13 — Concurrency patterns (heart of it)
- **Worker pool:** N goroutines pulling from a jobs channel, results on another.
- **Pipeline:** stage1 → stage2 → stage3 linked by channels.
- **Fan-in:** merge N channels into one.
- **Done channel:** propagate cancellation to a running goroutine.
- **Deliverable:** `scratch/day13_patterns/` with all four in separate files.

### Day 14 — context cancellation + goroutine leaks
- Wire `context.Context` into your worker pool so `cancel()` stops all workers.
- Introduce a **goroutine leak** intentionally. Detect with `runtime.NumGoroutine()`. Fix it.
- Read: [Go Concurrency Patterns: Context](https://go.dev/blog/context).
- **Deliverable:** a "how to detect goroutine leaks" note.

**End of Phase 2 review (Day 14):** you now know why every subsequent WebSocket line exists.

---

## Phase 3 — WebSocket chat core (Days 15–22)

**Why now:** everything above will be *applied* here. The chat is your concurrency final exam.

### Day 15 — WebSocket handshake
- Add `nhooyr.io/websocket` (or `gorilla/websocket` — pick, don't waffle).
- New endpoint `GET /api/v1/ws` behind auth middleware.
- Echo server: accept connection, echo any message back.
- **Deliverable:** working echo WS with auth.

### Day 16 — Client abstraction
- One `Client` struct per connection: `readPump` goroutine, `writePump` goroutine.
- `send chan []byte` buffered channel between them.
- Handle write deadlines + close on error.
- **Deliverable:** `internal/chat/client.go`.

### Day 17 — Hub pattern
- Central `Hub` goroutine with `register / unregister / broadcast` channels.
- Set of active clients (map, guarded correctly).
- Broadcast to all clients.
- **Deliverable:** `internal/chat/hub.go`. Two browser tabs, send a message, both receive.

### Day 18 — Rooms
- `rooms` and `room_members` tables + migrations.
- `POST /rooms`, `POST /rooms/:id/join`, `GET /rooms`.
- Hub tracks clients-per-room, broadcasts are room-scoped.
- **Deliverable:** rooms end-to-end. Auth enforced.

### Day 19 — Message persistence
- `messages` table.
- Save incoming message *before* broadcasting.
- `GET /rooms/:id/messages?before=&limit=` for history.
- **Deliverable:** history endpoint + integration test.

### Day 20 — Ping/pong + slow-client backpressure
- Server ping every 30s, disconnect on missed pong.
- If a client's `send` buffer is full, drop them (don't block the hub).
- Read: [gorilla/websocket chat example](https://github.com/gorilla/websocket/tree/main/examples/chat).
- **Deliverable:** chaos test — a slow client can't stall the hub.

### Day 21 — Presence
- Redis `SET user:{id}:online 1 EX 30` (introduces Redis to the stack).
- `GET /rooms/:id/presence` returns online users.
- Alternative if skipping Redis this day: in-memory presence map on the hub (still learn the pattern).
- **Deliverable:** presence visible in UI (curl for now).

### Day 22 — Load test locally
- Use `k6` or `hey` or a Go script to open 500 concurrent WS connections.
- Watch `runtime.NumGoroutine()`, memory. Fix anything that leaks or panics.
- **Deliverable:** notes on max concurrent connections on your machine.

**End of Phase 3:** you have a real chat app running on one instance.

---

## Phase 4 — Production infrastructure (Days 23–38)

This is what turns "a Go program" into "a service."

### Days 23–24 — Docker + docker-compose
- Multi-stage Dockerfile (build stage + `distroless` runtime).
- `docker-compose.yml` with app + MySQL + Redis.
- `make up` / `make down`.
- **Deliverable:** `docker compose up` gives you the whole stack.

### Days 25–26 — CI (GitHub Actions)
- Workflow: on PR → `go vet`, `golangci-lint`, `go test -race`, `go build`.
- On main → build + push a Docker image to GHCR.
- Cache Go modules between runs.
- **Deliverable:** green CI badge in README.

### Days 27–28 — sqlc adoption
- Replace one hand-written repo method with sqlc-generated code (start with `FindByUsername`).
- Understand: SQL as source of truth, Go types generated.
- Migrate the rest incrementally.
- **Deliverable:** whole `user` repo backed by sqlc.

### Days 29–30 — testcontainers integration tests
- Spin up a real MySQL container from tests.
- Rewrite integration tests to use it (no more mocking DB).
- Run in CI too.
- **Deliverable:** `make test-integration` uses testcontainers, green in CI.

### Days 31–32 — Observability: metrics
- Add Prometheus client library.
- Export request count, latency histogram, active WS connections, goroutine count.
- Add `/metrics` endpoint (unauth or private).
- Local Grafana via compose.
- **Deliverable:** you can see a live dashboard.

### Days 33–34 — Observability: tracing
- OpenTelemetry SDK.
- Trace: incoming HTTP → service → repo → DB call. Each span with attributes.
- Export to a local Jaeger container.
- **Deliverable:** flame graph of a login request.

### Day 35 — Observability: structured logs shipping
- Ensure `slog` handler emits JSON.
- Add trace ID + request ID to every log line.
- Compose stack: Loki + Grafana Explore for logs.
- **Deliverable:** grep-able logs correlated with traces.

### Days 36–37 — Redis for real scaling
- Redis pub/sub for cross-instance WebSocket broadcast.
- Redis-backed rate limiter (replaces in-memory).
- Redis for JWT revocation list.
- Run 2 instances behind nginx compose, prove messages fan out.
- **Deliverable:** 2 instances, chat works across them.

### Day 38 — Secrets + config hardening
- Kill `.env` in the "prod path." Load from environment only.
- Add a config validator that panics at startup on missing/invalid config.
- Distinguish `/livez` (I'm running) vs `/readyz` (I'm ready for traffic).
- **Deliverable:** shutdown flips readyz to 503 before draining.

**End of Phase 4:** you now have a service that could survive a real deployment.

---

## Phase 5 — Deploy, load-test, harden (Days 39–45)

### Days 39–40 — Deploy to a real cloud
- Pick one: Fly.io (easiest for Go), Railway, or a small GKE / EKS cluster if you want to learn Kubernetes.
- Managed MySQL (PlanetScale, Aiven, or cloud-native).
- Managed Redis (Upstash / cloud-native).
- Real domain, TLS via caddy/nginx or the platform's cert manager.
- **Deliverable:** public URL, `curl` it from your phone.

### Days 41–42 — Serious load test
- k6 script: N users, ramp up, mixed HTTP + WS traffic.
- Find your breaking point. Read metrics + traces to know why.
- Fix the first bottleneck (usually a DB query missing an index).
- **Deliverable:** load-test report with graphs.

### Day 43 — Security pass
- Add `helmet`-style headers via middleware.
- Add CSP if you have a frontend.
- Rotate JWT signing key procedure (`kid` header + key store).
- Run `gosec` and `govulncheck`.
- **Deliverable:** a security findings + fixes report.

### Day 44 — Documentation pass
- Update README: architecture diagram, how to run, how to deploy.
- OpenAPI spec generated (`oapi-codegen` or `swag`).
- ADRs (architecture decision records) for 3 key choices you made.
- **Deliverable:** anyone can clone + run in 10 min.

### Day 45 — Portfolio-ready polish + retrospective
- Record a 5-min video walkthrough of the app + architecture.
- Post the repo publicly. Write a blog post ("what I learned building a Go chat app in 45 days").
- Personal retro: what took longer than expected, what stuck, what didn't.
- **Deliverable:** shareable link + retro doc.

---

## Weekly review ritual (every Sunday)

1. Re-open `learning-checklist.html`. Tick what you covered.
2. Re-read the week's `notes-*.md`. Any item that made you go "wait, why?" → add to next week.
3. Skim closed items — can you still explain them to a rubber duck in 2 minutes? If not, revisit.
4. Update this file's "current day" pointer:

    - [ ] Current day: Day __
    - [ ] Blocked on: __
    - [ ] Confident on: __
    - [ ] Wobbly on: __

---

## Resources (bookmark these, don't drown in them)

- **The Go Programming Language** by Donovan & Kernighan — reference for language depth
- **100 Go Mistakes and How to Avoid Them** by Teiva Harsanyi — reads like a checklist itself
- **Effective Go** — https://go.dev/doc/effective_go
- **Go Code Review Comments** — https://github.com/golang/go/wiki/CodeReviewComments
- **Ardan Labs blog** (Bill Kennedy) — best on Go internals & memory
- **Dave Cheney's blog** — go-to for practical wisdom
- **Uber Go Style Guide** — https://github.com/uber-go/guide
- **Go Concurrency Patterns talk** (Rob Pike) — https://www.youtube.com/watch?v=f6kdp27TYZs
- **Gophercon talks on YouTube** — pick topics as they come up

Don't read all of these front to back. Use them as reference when a topic hits.

---

## Non-negotiables

- No skipping the notes.
- No skipping the tests.
- No pushing to main without CI green (once CI exists).
- No "I'll come back to that gotcha later" — write it down NOW in notes.

---

## Progress log

| Date | Day | Delivered | Blockers | Next |
|------|-----|-----------|----------|------|
| 2026-07-14 | 0 | Plan committed | — | Day 1: types |
| | | | | |
