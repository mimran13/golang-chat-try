# Learning Tasks - Codebase Expansion

A collection of tasks to help you learn Go by expanding this task manager API. Tasks are organized by difficulty level.

---

## Beginner Tasks

### 1. Add Task Tags/Labels

**Goal:** Allow users to categorize tasks with tags.

**Requirements:**
- Create a `Tag` model with `id`, `name`, `color` (hex), `user_id`
- Implement many-to-many relationship between Task and Tag
- Create CRUD endpoints:
  - `POST /api/v1/tags` - create tag
  - `GET /api/v1/tags` - list user's tags
  - `PUT /api/v1/tags/{id}` - update tag
  - `DELETE /api/v1/tags/{id}` - delete tag
- Add `tag_ids` field to create/update task requests
- Add `tags` array to task responses
- Allow filtering tasks by tag: `GET /api/v1/tasks?tag=work`

**Concepts to Learn:**
- GORM many-to-many relationships
- Junction tables (`task_tags`)
- Preloading associations

**Files to Create/Modify:**
- `internal/models/tag.go` (new)
- `internal/repository/tag_repository.go` (new)
- `internal/services/tag_service.go` (new)
- `internal/handlers/tag_handler.go` (new)
- `internal/dto/tag.go` (new)
- `internal/models/task.go` (add Tags field)
- `internal/repository/task_repository.go` (preload tags)
- `cmd/api/main.go` (wire up new routes)

---

### 2. Add Password Change Endpoint

**Goal:** Allow authenticated users to change their password.

**Requirements:**
- Create endpoint: `PUT /api/v1/auth/password`
- Request body:
  ```json
  {
    "current_password": "oldpass123",
    "new_password": "newpass456",
    "confirm_password": "newpass456"
  }
  ```
- Validate current password matches stored hash
- Validate new password meets requirements (min 8 chars)
- Validate new_password matches confirm_password
- Return appropriate errors for each validation failure

**Concepts to Learn:**
- Password verification with bcrypt
- Custom validation logic
- Meaningful error messages

**Files to Create/Modify:**
- `internal/dto/auth.go` (add ChangePasswordRequest)
- `internal/services/auth_service.go` (add ChangePassword method)
- `internal/handlers/auth_handler.go` (add ChangePassword handler)
- `cmd/api/main.go` (add route)

---

### 3. Add User Profile Update

**Goal:** Allow users to update their profile information.

**Requirements:**
- Create endpoint: `PUT /api/v1/auth/profile`
- Allow updating: `username`, `email`
- Handle unique constraint violations gracefully:
  - "Username already taken"
  - "Email already in use"
- Return updated user profile

**Concepts to Learn:**
- Partial updates (only update provided fields)
- Database constraint error handling
- GORM error type checking

**Files to Create/Modify:**
- `internal/dto/auth.go` (add UpdateProfileRequest)
- `internal/services/auth_service.go` (add UpdateProfile method)
- `internal/handlers/auth_handler.go` (add UpdateProfile handler)
- `internal/repository/user_repository.go` (add Update method)

---

## Intermediate Tasks

### 4. Add Task Comments

**Goal:** Allow users to comment on tasks.

**Requirements:**
- Create `Comment` model:
  - `id`, `task_id`, `user_id`, `content`, `created_at`, `updated_at`
- Users can only comment on their own tasks
- Endpoints:
  - `POST /api/v1/tasks/{id}/comments` - add comment
  - `GET /api/v1/tasks/{id}/comments` - list comments (paginated)
  - `DELETE /api/v1/tasks/{id}/comments/{comment_id}` - delete own comment
- Include `comment_count` in task list/detail responses
- Order comments by created_at descending (newest first)

**Concepts to Learn:**
- Nested RESTful resources
- Belongs-to relationships (Comment belongs to Task and User)
- Aggregation queries (COUNT)
- Preloading with conditions

**Files to Create/Modify:**
- `internal/models/comment.go` (new)
- `internal/repository/comment_repository.go` (new)
- `internal/services/comment_service.go` (new)
- `internal/handlers/comment_handler.go` (new)
- `internal/dto/comment.go` (new)
- `internal/dto/task.go` (add comment_count to response)

---

### 5. Add Subtasks/Checklist

**Goal:** Allow tasks to have a checklist of subtasks.

**Requirements:**
- Create `Subtask` model:
  - `id`, `task_id`, `title`, `completed`, `position` (for ordering)
- Endpoints:
  - `POST /api/v1/tasks/{id}/subtasks` - add subtask
  - `GET /api/v1/tasks/{id}/subtasks` - list subtasks
  - `PUT /api/v1/tasks/{id}/subtasks/{subtask_id}` - update subtask
  - `PATCH /api/v1/tasks/{id}/subtasks/{subtask_id}/toggle` - toggle completion
  - `DELETE /api/v1/tasks/{id}/subtasks/{subtask_id}` - delete subtask
  - `PUT /api/v1/tasks/{id}/subtasks/reorder` - reorder subtasks
- Add to task response:
  - `subtasks` array
  - `subtask_count` (total)
  - `completed_subtask_count`
  - `progress` (percentage, computed)

**Concepts to Learn:**
- One-to-many relationships
- Computed/virtual fields
- Array position management
- Batch updates

**Files to Create/Modify:**
- `internal/models/subtask.go` (new)
- `internal/repository/subtask_repository.go` (new)
- `internal/services/subtask_service.go` (new)
- `internal/handlers/subtask_handler.go` (new)
- `internal/dto/subtask.go` (new)
- `internal/dto/task.go` (add subtask fields)

---

### 6. Add Task Activity/Audit Log

**Goal:** Track all changes made to tasks.

**Requirements:**
- Create `TaskActivity` model:
  - `id`, `task_id`, `user_id`, `action` (created, updated, status_changed, deleted, restored)
  - `old_value`, `new_value` (JSON fields for storing changes)
  - `created_at`
- Automatically log activity using GORM hooks
- Endpoint: `GET /api/v1/tasks/{id}/activity` (paginated)
- Activity types to track:
  - Task created
  - Task updated (which fields changed)
  - Status changed (from → to)
  - Task deleted
  - Task restored

**Concepts to Learn:**
- GORM hooks (AfterCreate, AfterUpdate, AfterDelete)
- JSON columns in database
- Event sourcing basics
- Diffing old vs new values

**Files to Create/Modify:**
- `internal/models/task_activity.go` (new)
- `internal/repository/task_activity_repository.go` (new)
- `internal/services/task_activity_service.go` (new)
- `internal/handlers/task_activity_handler.go` (new)
- `internal/dto/task_activity.go` (new)
- `internal/models/task.go` (add hooks)

---

### 7. Implement Soft Delete Recovery

**Goal:** Allow users to recover deleted tasks.

**Requirements:**
- Endpoint to list deleted tasks: `GET /api/v1/tasks/deleted`
  - Should support pagination
  - Only show tasks deleted in last 30 days
- Endpoint to restore: `POST /api/v1/tasks/{id}/restore`
- Endpoint to permanently delete: `DELETE /api/v1/tasks/{id}/permanent`
- Add `deleted_at` to task responses when listing deleted tasks

**Concepts to Learn:**
- GORM Unscoped() for querying soft-deleted records
- Restore functionality (set deleted_at to NULL)
- Hard delete vs soft delete
- Time-based filtering

**Files to Create/Modify:**
- `internal/repository/task_repository.go` (add methods)
- `internal/services/task_service.go` (add methods)
- `internal/handlers/task_handler.go` (add handlers)
- `internal/dto/task.go` (add DeletedTaskResponse)

---

## Advanced Tasks

### 8. Add Redis Caching

**Goal:** Implement caching to improve read performance.

**Requirements:**
- Set up Redis connection (add to docker-compose)
- Cache task list queries with key pattern: `user:{id}:tasks:{hash_of_filters}`
- Cache individual task: `task:{id}`
- Cache TTL: 5 minutes
- Invalidate cache on:
  - Task created → invalidate list cache
  - Task updated → invalidate task cache + list cache
  - Task deleted → invalidate task cache + list cache
- Add cache hit/miss to response headers (X-Cache: HIT/MISS)

**Concepts to Learn:**
- go-redis client
- Cache-aside pattern
- Cache key design
- Cache invalidation strategies
- JSON serialization for cache

**Files to Create/Modify:**
- `docker-compose.yml` (add Redis service)
- `internal/config/config.go` (add Redis config)
- `internal/cache/redis.go` (new - connection setup)
- `internal/cache/task_cache.go` (new - task caching logic)
- `internal/services/task_service.go` (integrate caching)
- `.env.example` (add Redis config)

---

### 9. Add WebSocket for Real-time Updates

**Goal:** Push real-time updates to connected clients.

**Requirements:**
- WebSocket endpoint: `GET /ws`
- Authenticate WebSocket connections via query param token
- Send events when:
  - Task created: `{ "type": "task.created", "data": {...} }`
  - Task updated: `{ "type": "task.updated", "data": {...} }`
  - Task deleted: `{ "type": "task.deleted", "data": { "id": 123 } }`
- Only send events for the connected user's tasks
- Handle connection lifecycle (connect, disconnect, reconnect)
- Implement heartbeat/ping-pong for connection health

**Concepts to Learn:**
- gorilla/websocket library
- Connection management (hub pattern)
- Concurrent map access with sync.RWMutex
- Channel-based event broadcasting
- Graceful connection handling

**Files to Create/Modify:**
- `internal/websocket/hub.go` (new - connection manager)
- `internal/websocket/client.go` (new - client connection)
- `internal/websocket/handler.go` (new - upgrade handler)
- `internal/services/task_service.go` (broadcast events)
- `cmd/api/main.go` (add WebSocket route)
- `go.mod` (add gorilla/websocket)

---

### 10. Add Email Verification Flow

**Goal:** Verify user email addresses before allowing login.

**Requirements:**
- Add fields to User model: `email_verified` (bool), `verification_token`, `verification_sent_at`
- On registration:
  - Generate secure verification token
  - Send verification email (extend notification service)
  - Return message: "Please verify your email"
- Verification endpoint: `GET /api/v1/auth/verify?token=xxx`
  - Validate token and mark email as verified
  - Token expires after 24 hours
- Resend verification: `POST /api/v1/auth/resend-verification`
  - Rate limit: 1 per 5 minutes
- Block login for unverified users with clear error message
- Optional: Add verified badge to user profile

**Concepts to Learn:**
- Secure token generation (crypto/rand)
- Token expiration handling
- Email templating
- Rate limiting per action
- State management in authentication flow

**Files to Create/Modify:**
- `internal/models/user.go` (add fields)
- `internal/services/auth_service.go` (verification logic)
- `internal/services/notification_service.go` (email sending)
- `internal/handlers/auth_handler.go` (new endpoints)
- `internal/dto/auth.go` (new request/response types)
- Database migration for new fields

---

### 11. Add Role-Based Access Control (RBAC)

**Goal:** Implement user roles and permissions.

**Requirements:**
- Add `role` field to User model: `admin`, `user` (default)
- Create authorization middleware that checks roles
- Admin-only endpoints:
  - `GET /api/v1/admin/users` - list all users
  - `GET /api/v1/admin/users/{id}` - get user details
  - `PUT /api/v1/admin/users/{id}/role` - change user role
  - `DELETE /api/v1/admin/users/{id}` - delete user
  - `GET /api/v1/admin/stats` - system-wide statistics
- Regular users cannot access admin routes (403 Forbidden)
- Include role in JWT claims
- Include role in user profile response

**Concepts to Learn:**
- Role-based authorization
- Middleware chaining
- JWT claims extension
- Admin vs user separation
- Permission checking patterns

**Files to Create/Modify:**
- `internal/models/user.go` (add Role field)
- `internal/middleware/authorization.go` (new)
- `internal/handlers/admin_handler.go` (new)
- `internal/services/admin_service.go` (new)
- `internal/dto/admin.go` (new)
- `internal/utils/jwt.go` (add role to claims)
- `cmd/api/main.go` (add admin routes)

---

---

# Go Core Concepts - Deep Dive Tasks

These tasks focus on mastering Go's core language features and patterns. Essential for senior/staff level interviews.

---

## Goroutines & Concurrency

### G1. Build a Worker Pool for Bulk Task Processing

**Goal:** Process multiple tasks concurrently with controlled parallelism.

**Requirements:**
- Create endpoint: `POST /api/v1/tasks/bulk/process`
- Accept array of task IDs to process (e.g., send reminders, generate reports)
- Implement worker pool pattern:
  - Configurable number of workers (default: 5)
  - Jobs channel to distribute work
  - Results channel to collect outcomes
  - WaitGroup to synchronize completion
- Return processing results: `{ "processed": 10, "failed": 2, "errors": [...] }`
- Add timeout for entire operation (context with deadline)

**Implementation Skeleton:**
```go
type Job struct {
    TaskID uint
    Action string
}

type Result struct {
    TaskID  uint
    Success bool
    Error   error
}

func (s *Service) ProcessBulk(ctx context.Context, jobs []Job, workers int) []Result {
    jobsChan := make(chan Job, len(jobs))
    resultsChan := make(chan Result, len(jobs))

    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go s.worker(ctx, &wg, jobsChan, resultsChan)
    }

    // Send jobs
    for _, job := range jobs {
        jobsChan <- job
    }
    close(jobsChan)

    // Wait and collect
    go func() {
        wg.Wait()
        close(resultsChan)
    }()

    var results []Result
    for result := range resultsChan {
        results = append(results, result)
    }
    return results
}
```

**Concepts to Learn:**
- Goroutine lifecycle management
- Buffered vs unbuffered channels
- sync.WaitGroup for coordination
- Channel closing semantics
- Fan-out/fan-in pattern

**Interview Questions This Prepares You For:**
- "How would you process 1 million records concurrently?"
- "Explain the worker pool pattern"
- "What happens if you don't close a channel?"

---

### G2. Implement Concurrent Task Search with Timeout

**Goal:** Search across multiple data sources concurrently with timeout.

**Requirements:**
- Create endpoint: `GET /api/v1/search?q=keyword`
- Search concurrently in:
  - Task titles
  - Task descriptions
  - Comments (if implemented)
  - Tags (if implemented)
- Use goroutines for parallel search
- Implement timeout (e.g., 3 seconds max)
- Return partial results if some searches timeout
- Use context for cancellation propagation

**Implementation Pattern:**
```go
func (s *Service) Search(ctx context.Context, query string) (*SearchResults, error) {
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    results := &SearchResults{}
    var mu sync.Mutex
    var wg sync.WaitGroup

    // Search titles
    wg.Add(1)
    go func() {
        defer wg.Done()
        if r, err := s.searchTitles(ctx, query); err == nil {
            mu.Lock()
            results.Titles = r
            mu.Unlock()
        }
    }()

    // Search descriptions... (similar)

    wg.Wait()
    return results, ctx.Err()
}
```

**Concepts to Learn:**
- context.WithTimeout / context.WithCancel
- Context propagation through call stack
- Mutex for protecting shared state
- Graceful degradation with partial results
- Error handling in concurrent code

---

## Channels & Select

### C1. Build a Real-time Task Event Stream

**Goal:** Stream task events to multiple subscribers using channels.

**Requirements:**
- Create an EventBroker that manages subscribers
- Events: task.created, task.updated, task.deleted, task.status_changed
- Multiple subscribers can listen to same events
- Subscribers can filter by event type
- Non-blocking publish (don't block if subscriber is slow)
- Automatic cleanup of dead subscribers

**Implementation:**
```go
type Event struct {
    Type      string
    TaskID    uint
    UserID    uint
    Payload   interface{}
    Timestamp time.Time
}

type EventBroker struct {
    subscribers map[chan Event][]string // channel -> event types
    mu          sync.RWMutex
}

func (b *EventBroker) Subscribe(eventTypes ...string) <-chan Event {
    ch := make(chan Event, 100) // buffered to prevent blocking
    b.mu.Lock()
    b.subscribers[ch] = eventTypes
    b.mu.Unlock()
    return ch
}

func (b *EventBroker) Publish(event Event) {
    b.mu.RLock()
    defer b.mu.RUnlock()

    for ch, types := range b.subscribers {
        if b.shouldReceive(types, event.Type) {
            select {
            case ch <- event:
            default:
                // Subscriber too slow, skip (or log warning)
            }
        }
    }
}
```

**Concepts to Learn:**
- Channel directions (`<-chan` vs `chan<-`)
- Buffered channels for async communication
- Select with default for non-blocking operations
- Pub/sub pattern in Go
- Thread-safe subscriber management

---

### C2. Implement Rate Limiter Using Ticker and Channels

**Goal:** Build a token bucket rate limiter from scratch.

**Requirements:**
- Create `RateLimiter` struct with configurable rate and burst
- Use `time.Ticker` to refill tokens
- `Allow() bool` - check if request is allowed
- `Wait(ctx context.Context) error` - block until allowed or context cancelled
- Must be thread-safe for concurrent use

**Implementation:**
```go
type RateLimiter struct {
    tokens     chan struct{}
    ticker     *time.Ticker
    maxTokens  int
    done       chan struct{}
}

func NewRateLimiter(rate time.Duration, burst int) *RateLimiter {
    rl := &RateLimiter{
        tokens:    make(chan struct{}, burst),
        ticker:    time.NewTicker(rate),
        maxTokens: burst,
        done:      make(chan struct{}),
    }

    // Fill initial tokens
    for i := 0; i < burst; i++ {
        rl.tokens <- struct{}{}
    }

    // Start refill goroutine
    go rl.refill()
    return rl
}

func (rl *RateLimiter) refill() {
    for {
        select {
        case <-rl.ticker.C:
            select {
            case rl.tokens <- struct{}{}:
            default: // bucket full
            }
        case <-rl.done:
            return
        }
    }
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
    select {
    case <-rl.tokens:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

**Concepts to Learn:**
- time.Ticker for periodic operations
- Select with multiple channels
- Channel as semaphore (buffered channel for counting)
- Graceful shutdown of background goroutines
- Empty struct `struct{}{}` for signaling (zero memory)

---

### C3. Build a Request Coalescer (Singleflight Pattern)

**Goal:** Deduplicate concurrent identical requests.

**Scenario:** Multiple users request the same expensive computation (e.g., stats). Instead of running it 10 times, run once and share result.

**Requirements:**
- Create `Coalescer` that groups identical in-flight requests
- Key-based grouping (same key = same request)
- First request triggers computation
- Subsequent requests wait for first to complete
- All waiters receive same result
- Cleanup after completion

**Implementation:**
```go
type call struct {
    wg     sync.WaitGroup
    result interface{}
    err    error
}

type Coalescer struct {
    mu    sync.Mutex
    calls map[string]*call
}

func (c *Coalescer) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
    c.mu.Lock()
    if c.calls == nil {
        c.calls = make(map[string]*call)
    }

    if existing, ok := c.calls[key]; ok {
        c.mu.Unlock()
        existing.wg.Wait()
        return existing.result, existing.err
    }

    call := &call{}
    call.wg.Add(1)
    c.calls[key] = call
    c.mu.Unlock()

    call.result, call.err = fn()
    call.wg.Done()

    c.mu.Lock()
    delete(c.calls, key)
    c.mu.Unlock()

    return call.result, call.err
}
```

**Use Case:** `GET /api/v1/tasks/stats` - if 10 requests come simultaneously, compute only once.

**Concepts to Learn:**
- sync.WaitGroup for waiting on completion
- Map as in-flight request tracker
- Mutex for map access
- Shared result pattern
- This is exactly what `golang.org/x/sync/singleflight` does

---

## Context Package

### CTX1. Implement Request Tracing with Context

**Goal:** Propagate request ID and tracing info through entire request lifecycle.

**Requirements:**
- Middleware that generates/extracts request ID
- Store in context: request ID, user ID, start time
- All logs include request ID automatically
- Pass context through: Handler → Service → Repository
- Add request duration to response headers

**Implementation:**
```go
type contextKey string

const (
    RequestIDKey contextKey = "request_id"
    UserIDKey    contextKey = "user_id"
    StartTimeKey contextKey = "start_time"
)

func TracingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }

        ctx := context.WithValue(c.Request.Context(), RequestIDKey, requestID)
        ctx = context.WithValue(ctx, StartTimeKey, time.Now())

        c.Request = c.Request.WithContext(ctx)
        c.Header("X-Request-ID", requestID)

        c.Next()

        duration := time.Since(ctx.Value(StartTimeKey).(time.Time))
        c.Header("X-Response-Time", duration.String())
    }
}

// In service layer
func (s *TaskService) GetTask(ctx context.Context, id uint) (*Task, error) {
    requestID := ctx.Value(RequestIDKey).(string)
    s.logger.Info("fetching task", zap.String("request_id", requestID), zap.Uint("task_id", id))
    // ...
}
```

**Concepts to Learn:**
- context.WithValue for request-scoped data
- Custom context key types (avoid collisions)
- Context propagation patterns
- Type assertions from context values
- Why context should be first parameter

---

### CTX2. Implement Graceful Shutdown with Context

**Goal:** Properly shutdown all components when server stops.

**Requirements:**
- Listen for OS signals (SIGINT, SIGTERM)
- Create cancellable context for application lifetime
- Propagate cancellation to:
  - HTTP server (stop accepting new requests)
  - Background workers (finish current job, exit)
  - Database connections (close pool)
  - WebSocket connections (send close frame)
- Wait for all components with timeout
- Force exit if graceful shutdown takes too long

**Implementation:**
```go
func main() {
    // Root context for entire application
    ctx, cancel := context.WithCancel(context.Background())

    // Signal handling
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sigChan
        log.Println("Shutdown signal received")
        cancel()
    }()

    // Start components with context
    server := startServer(ctx)
    workers := startWorkers(ctx)

    // Wait for shutdown
    <-ctx.Done()

    // Graceful shutdown with timeout
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer shutdownCancel()

    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        server.Shutdown(shutdownCtx)
    }()

    go func() {
        defer wg.Done()
        workers.Stop(shutdownCtx)
    }()

    wg.Wait()
    log.Println("Graceful shutdown complete")
}
```

**Concepts to Learn:**
- Signal handling in Go
- Context cancellation propagation
- Coordinated shutdown of multiple components
- Timeout for shutdown operations
- Clean resource cleanup

---

## Interfaces & Composition

### I1. Implement Repository Pattern with Interfaces

**Goal:** Make repositories fully testable with interface-based design.

**Requirements:**
- Define interfaces for all repositories
- Implement mock repositories for testing
- Use dependency injection via interfaces
- Create integration tests with real DB
- Create unit tests with mocks

**Implementation:**
```go
// interfaces.go
type TaskRepository interface {
    Create(ctx context.Context, task *Task) error
    GetByID(ctx context.Context, id uint) (*Task, error)
    Update(ctx context.Context, task *Task) error
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context, filter TaskFilter) ([]Task, int64, error)
}

// mock_task_repository.go
type MockTaskRepository struct {
    tasks  map[uint]*Task
    nextID uint
    mu     sync.RWMutex
}

func (m *MockTaskRepository) Create(ctx context.Context, task *Task) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.nextID++
    task.ID = m.nextID
    m.tasks[task.ID] = task
    return nil
}

// Service uses interface, not concrete type
type TaskService struct {
    repo TaskRepository // interface, not *GormTaskRepository
}
```

**Concepts to Learn:**
- Interface definition and implementation
- Implicit interface satisfaction (no `implements` keyword)
- Dependency injection in Go
- Interface segregation (small, focused interfaces)
- Testing with mock implementations

---

### I2. Build a Pluggable Storage Backend

**Goal:** Support multiple storage backends (MySQL, PostgreSQL, SQLite, In-Memory).

**Requirements:**
- Define `Storage` interface with all needed methods
- Implement for each backend
- Factory function to create storage based on config
- Storage can be swapped without changing business logic
- Add in-memory storage for testing

**Implementation:**
```go
type Storage interface {
    TaskRepository
    UserRepository
    // Add other repositories

    Close() error
    Ping(ctx context.Context) error
    Migrate() error
}

type MySQLStorage struct {
    db *gorm.DB
}

type PostgresStorage struct {
    db *gorm.DB
}

type InMemoryStorage struct {
    tasks map[uint]*Task
    users map[uint]*User
    mu    sync.RWMutex
}

func NewStorage(cfg *Config) (Storage, error) {
    switch cfg.DatabaseDriver {
    case "mysql":
        return NewMySQLStorage(cfg.DatabaseDSN)
    case "postgres":
        return NewPostgresStorage(cfg.DatabaseDSN)
    case "memory":
        return NewInMemoryStorage(), nil
    default:
        return nil, fmt.Errorf("unsupported driver: %s", cfg.DatabaseDriver)
    }
}
```

**Concepts to Learn:**
- Interface composition (embedding interfaces)
- Factory pattern in Go
- Strategy pattern via interfaces
- Configuration-driven implementations

---

### I3. Implement Middleware Chain with Interfaces

**Goal:** Build a composable middleware system from scratch.

**Requirements:**
- Define `Middleware` interface/type
- Implement middleware: logging, auth, rate limiting, recovery, tracing
- Support middleware ordering
- Allow per-route middleware configuration
- Implement `Use()`, `Group()` patterns

**Implementation:**
```go
type Handler func(ctx *Context) error

type Middleware func(Handler) Handler

// Chain multiple middleware
func Chain(middlewares ...Middleware) Middleware {
    return func(final Handler) Handler {
        for i := len(middlewares) - 1; i >= 0; i-- {
            final = middlewares[i](final)
        }
        return final
    }
}

// Example middleware
func LoggingMiddleware(next Handler) Handler {
    return func(ctx *Context) error {
        start := time.Now()
        err := next(ctx)
        log.Printf("%s %s %v", ctx.Method, ctx.Path, time.Since(start))
        return err
    }
}

func AuthMiddleware(next Handler) Handler {
    return func(ctx *Context) error {
        token := ctx.GetHeader("Authorization")
        if !validateToken(token) {
            return ErrUnauthorized
        }
        return next(ctx)
    }
}

// Usage
handler := Chain(
    LoggingMiddleware,
    AuthMiddleware,
    RateLimitMiddleware,
)(finalHandler)
```

**Concepts to Learn:**
- Function types as interfaces
- Higher-order functions
- Decorator pattern in Go
- Middleware composition
- This is how Gin/Echo middleware works internally

---

## Error Handling

### E1. Implement Comprehensive Error Handling System

**Goal:** Build production-grade error handling with wrapping, types, and stack traces.

**Requirements:**
- Custom error types with codes, messages, HTTP status
- Error wrapping with context (`fmt.Errorf` with `%w`)
- Error unwrapping and type checking (`errors.Is`, `errors.As`)
- Stack trace capture for debugging
- Consistent error response format
- Logging with full context

**Implementation:**
```go
type ErrorCode string

const (
    ErrCodeNotFound       ErrorCode = "NOT_FOUND"
    ErrCodeValidation     ErrorCode = "VALIDATION_ERROR"
    ErrCodeUnauthorized   ErrorCode = "UNAUTHORIZED"
    ErrCodeInternal       ErrorCode = "INTERNAL_ERROR"
    ErrCodeConflict       ErrorCode = "CONFLICT"
)

type AppError struct {
    Code       ErrorCode
    Message    string
    HTTPStatus int
    Err        error  // wrapped error
    Stack      string // stack trace
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

func (e *AppError) Unwrap() error {
    return e.Err
}

// Constructors
func NewNotFoundError(resource string, id interface{}) *AppError {
    return &AppError{
        Code:       ErrCodeNotFound,
        Message:    fmt.Sprintf("%s with ID %v not found", resource, id),
        HTTPStatus: http.StatusNotFound,
        Stack:      captureStack(),
    }
}

func WrapError(err error, message string) *AppError {
    return &AppError{
        Code:       ErrCodeInternal,
        Message:    message,
        HTTPStatus: http.StatusInternalServerError,
        Err:        err,
        Stack:      captureStack(),
    }
}

// Usage in service
func (s *TaskService) GetTask(ctx context.Context, id uint) (*Task, error) {
    task, err := s.repo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, NewNotFoundError("Task", id)
        }
        return nil, WrapError(err, "failed to fetch task")
    }
    return task, nil
}

// Error handler middleware
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            var appErr *AppError
            if errors.As(err, &appErr) {
                c.JSON(appErr.HTTPStatus, ErrorResponse{
                    Code:    string(appErr.Code),
                    Message: appErr.Message,
                })
                return
            }
            // Unknown error
            c.JSON(500, ErrorResponse{Code: "INTERNAL_ERROR", Message: "An error occurred"})
        }
    }
}
```

**Concepts to Learn:**
- Custom error types
- Error wrapping with `%w`
- `errors.Is` and `errors.As`
- Implementing `Error()` and `Unwrap()` interfaces
- Stack trace capture (`runtime.Stack`)
- Sentinel errors vs error types

---

### E2. Implement Panic Recovery with Detailed Reporting

**Goal:** Catch panics, log them with context, and return proper error response.

**Requirements:**
- Recover from panics in HTTP handlers
- Capture stack trace
- Log: request details, user info, stack trace
- Return 500 with request ID for debugging
- Optional: Send to error tracking service (Sentry pattern)

**Implementation:**
```go
func RecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if r := recover(); r != nil {
                // Capture stack
                stack := make([]byte, 4096)
                n := runtime.Stack(stack, false)

                // Get request context
                requestID := c.GetString("request_id")
                userID, _ := c.Get("user_id")

                // Log with full context
                logger.Error("panic recovered",
                    zap.Any("panic", r),
                    zap.String("request_id", requestID),
                    zap.Any("user_id", userID),
                    zap.String("method", c.Request.Method),
                    zap.String("path", c.Request.URL.Path),
                    zap.ByteString("stack", stack[:n]),
                )

                // Return error response
                c.AbortWithStatusJSON(500, gin.H{
                    "error":      "Internal server error",
                    "request_id": requestID,
                })
            }
        }()
        c.Next()
    }
}
```

**Concepts to Learn:**
- `recover()` function and when it works
- `defer` execution order
- `runtime.Stack` for stack traces
- Panic vs error - when to use each
- Never panic in libraries, only in main

---

## Sync Package Deep Dive

### S1. Implement Thread-Safe Cache with sync.RWMutex

**Goal:** Build an in-memory cache with proper concurrent access.

**Requirements:**
- Get, Set, Delete operations
- TTL support for entries
- Background cleanup of expired entries
- Use RWMutex (multiple readers, single writer)
- Benchmarks comparing Mutex vs RWMutex

**Implementation:**
```go
type CacheEntry struct {
    Value      interface{}
    Expiration time.Time
}

type Cache struct {
    data map[string]CacheEntry
    mu   sync.RWMutex
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    entry, ok := c.data[key]
    if !ok || time.Now().After(entry.Expiration) {
        return nil, false
    }
    return entry.Value, true
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.data[key] = CacheEntry{
        Value:      value,
        Expiration: time.Now().Add(ttl),
    }
}

// Background cleanup
func (c *Cache) StartCleanup(interval time.Duration, done <-chan struct{}) {
    ticker := time.NewTicker(interval)
    go func() {
        for {
            select {
            case <-ticker.C:
                c.cleanup()
            case <-done:
                ticker.Stop()
                return
            }
        }
    }()
}

func (c *Cache) cleanup() {
    c.mu.Lock()
    defer c.mu.Unlock()

    now := time.Now()
    for key, entry := range c.data {
        if now.After(entry.Expiration) {
            delete(c.data, key)
        }
    }
}
```

**Concepts to Learn:**
- sync.Mutex vs sync.RWMutex
- Lock contention and performance
- Background goroutines for maintenance
- Graceful shutdown of background tasks

---

### S2. Implement Connection Pool with sync.Pool

**Goal:** Reuse expensive objects to reduce GC pressure.

**Requirements:**
- Pool for reusable objects (buffers, connections, etc.)
- Create `BufferPool` for JSON encoding/decoding
- Benchmark memory allocation with/without pool
- Proper cleanup when returning to pool

**Implementation:**
```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func GetBuffer() *bytes.Buffer {
    return bufferPool.Get().(*bytes.Buffer)
}

func PutBuffer(buf *bytes.Buffer) {
    buf.Reset() // Important: clean before returning
    bufferPool.Put(buf)
}

// Usage in handler
func (h *Handler) GetTasks(c *gin.Context) {
    buf := GetBuffer()
    defer PutBuffer(buf)

    encoder := json.NewEncoder(buf)
    encoder.Encode(tasks)

    c.Data(200, "application/json", buf.Bytes())
}

// Benchmark
func BenchmarkWithPool(b *testing.B) {
    for i := 0; i < b.N; i++ {
        buf := GetBuffer()
        buf.WriteString("test")
        PutBuffer(buf)
    }
}

func BenchmarkWithoutPool(b *testing.B) {
    for i := 0; i < b.N; i++ {
        buf := new(bytes.Buffer)
        buf.WriteString("test")
    }
}
```

**Concepts to Learn:**
- sync.Pool usage and semantics
- Object lifecycle in pool
- GC pressure reduction
- When to use pools (high allocation rate)
- Benchmarking with `go test -bench`

---

### S3. Implement Once-Only Initialization with sync.Once

**Goal:** Ensure expensive initialization happens exactly once.

**Requirements:**
- Lazy initialization of database connection
- Config loading that happens once
- Thread-safe singleton pattern
- Handle initialization errors properly

**Implementation:**
```go
type Database struct {
    conn *sql.DB
    once sync.Once
    err  error
}

var db = &Database{}

func GetDB() (*sql.DB, error) {
    db.once.Do(func() {
        db.conn, db.err = sql.Open("mysql", os.Getenv("DATABASE_URL"))
        if db.err != nil {
            return
        }
        db.err = db.conn.Ping()
    })
    return db.conn, db.err
}

// With error retry pattern (sync.Once doesn't retry on error)
type LazyDB struct {
    conn *sql.DB
    mu   sync.Mutex
    done bool
}

func (l *LazyDB) Get() (*sql.DB, error) {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.done {
        return l.conn, nil
    }

    conn, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
    if err != nil {
        return nil, err // Will retry next call
    }

    l.conn = conn
    l.done = true
    return l.conn, nil
}
```

**Concepts to Learn:**
- sync.Once semantics
- Lazy initialization pattern
- Singleton in Go (and why to avoid it)
- Error handling with Once (it won't retry)

---

## Generics (Go 1.18+)

### GEN1. Build Generic Repository Layer

**Goal:** Create type-safe repository methods using generics.

**Requirements:**
- Generic CRUD operations
- Type-safe without interface{} casting
- Works with any model type
- Reduce code duplication

**Implementation:**
```go
type Model interface {
    GetID() uint
}

type Repository[T Model] struct {
    db *gorm.DB
}

func NewRepository[T Model](db *gorm.DB) *Repository[T] {
    return &Repository[T]{db: db}
}

func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
    return r.db.WithContext(ctx).Create(entity).Error
}

func (r *Repository[T]) GetByID(ctx context.Context, id uint) (*T, error) {
    var entity T
    err := r.db.WithContext(ctx).First(&entity, id).Error
    return &entity, err
}

func (r *Repository[T]) Update(ctx context.Context, entity *T) error {
    return r.db.WithContext(ctx).Save(entity).Error
}

func (r *Repository[T]) Delete(ctx context.Context, id uint) error {
    var entity T
    return r.db.WithContext(ctx).Delete(&entity, id).Error
}

func (r *Repository[T]) List(ctx context.Context, limit, offset int) ([]T, error) {
    var entities []T
    err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&entities).Error
    return entities, err
}

// Usage
taskRepo := NewRepository[Task](db)
userRepo := NewRepository[User](db)

task, err := taskRepo.GetByID(ctx, 1) // Returns *Task, not interface{}
```

**Concepts to Learn:**
- Generic type parameters `[T any]`
- Type constraints with interfaces
- Generic functions vs generic types
- When generics reduce vs increase complexity

---

### GEN2. Implement Generic Result Type (Rust-like)

**Goal:** Create a Result type for explicit error handling.

**Requirements:**
- `Result[T]` type that holds either value or error
- Methods: `IsOk()`, `IsErr()`, `Unwrap()`, `UnwrapOr(default)`
- `Map()` for transforming success values
- Use in service layer for cleaner error handling

**Implementation:**
```go
type Result[T any] struct {
    value T
    err   error
}

func Ok[T any](value T) Result[T] {
    return Result[T]{value: value}
}

func Err[T any](err error) Result[T] {
    return Result[T]{err: err}
}

func (r Result[T]) IsOk() bool {
    return r.err == nil
}

func (r Result[T]) IsErr() bool {
    return r.err != nil
}

func (r Result[T]) Unwrap() T {
    if r.err != nil {
        panic(r.err)
    }
    return r.value
}

func (r Result[T]) UnwrapOr(defaultValue T) T {
    if r.err != nil {
        return defaultValue
    }
    return r.value
}

func (r Result[T]) UnwrapOrElse(fn func() T) T {
    if r.err != nil {
        return fn()
    }
    return r.value
}

func Map[T, U any](r Result[T], fn func(T) U) Result[U] {
    if r.err != nil {
        return Err[U](r.err)
    }
    return Ok(fn(r.value))
}

// Usage
func (s *TaskService) GetTask(id uint) Result[*Task] {
    task, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return Err[*Task](err)
    }
    return Ok(task)
}

// Client code
result := taskService.GetTask(1)
if result.IsOk() {
    task := result.Unwrap()
}

// Or with default
task := taskService.GetTask(1).UnwrapOr(&Task{Title: "Default"})
```

**Concepts to Learn:**
- Generic type design
- Functional programming patterns in Go
- Optional/Result types
- Method chaining with generics

---

## Testing Patterns

### T1. Implement Table-Driven Tests with Subtests

**Goal:** Master Go's testing patterns.

**Requirements:**
- Table-driven tests for validation functions
- Subtests with `t.Run()`
- Parallel test execution
- Test helpers and cleanup
- Coverage for edge cases

**Implementation:**
```go
func TestValidateTask(t *testing.T) {
    tests := []struct {
        name    string
        task    CreateTaskRequest
        wantErr bool
        errField string
    }{
        {
            name:    "valid task",
            task:    CreateTaskRequest{Title: "Valid", Priority: "high"},
            wantErr: false,
        },
        {
            name:    "empty title",
            task:    CreateTaskRequest{Title: "", Priority: "high"},
            wantErr: true,
            errField: "title",
        },
        {
            name:    "title too long",
            task:    CreateTaskRequest{Title: strings.Repeat("a", 201), Priority: "high"},
            wantErr: true,
            errField: "title",
        },
        {
            name:    "invalid priority",
            task:    CreateTaskRequest{Title: "Valid", Priority: "invalid"},
            wantErr: true,
            errField: "priority",
        },
    }

    for _, tt := range tests {
        tt := tt // capture range variable
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel() // run subtests in parallel

            err := validate.Struct(tt.task)

            if tt.wantErr {
                if err == nil {
                    t.Errorf("expected error for field %s, got nil", tt.errField)
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
            }
        })
    }
}
```

**Concepts to Learn:**
- Table-driven test pattern
- `t.Run()` for subtests
- `t.Parallel()` for concurrent tests
- Range variable capture (`tt := tt`)
- Test naming conventions

---

### T2. Implement Integration Tests with Test Database

**Goal:** Test full request/response cycle with real database.

**Requirements:**
- Setup/teardown test database
- Test fixtures for consistent state
- HTTP handler testing with `httptest`
- Transaction rollback for test isolation
- Test configuration separate from prod

**Implementation:**
```go
type TestSuite struct {
    db      *gorm.DB
    router  *gin.Engine
    cleanup func()
}

func SetupTestSuite(t *testing.T) *TestSuite {
    // Use SQLite for tests
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)

    // Run migrations
    db.AutoMigrate(&User{}, &Task{})

    // Setup router with test DB
    router := setupRouter(db)

    return &TestSuite{
        db:     db,
        router: router,
        cleanup: func() {
            sqlDB, _ := db.DB()
            sqlDB.Close()
        },
    }
}

func (s *TestSuite) CreateTestUser(t *testing.T) *User {
    user := &User{
        Username: "testuser",
        Email:    "test@example.com",
    }
    user.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
    require.NoError(t, s.db.Create(user).Error)
    return user
}

func TestCreateTask(t *testing.T) {
    suite := SetupTestSuite(t)
    defer suite.cleanup()

    user := suite.CreateTestUser(t)
    token := generateTestToken(user)

    body := `{"title": "Test Task", "priority": "high"}`
    req := httptest.NewRequest("POST", "/api/v1/tasks", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    w := httptest.NewRecorder()
    suite.router.ServeHTTP(w, req)

    assert.Equal(t, 201, w.Code)

    var response TaskResponse
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, "Test Task", response.Title)
}
```

**Concepts to Learn:**
- httptest.NewRecorder for HTTP testing
- Test fixtures and helpers
- In-memory database for tests
- Test isolation strategies
- Setup/teardown patterns

---

### T3. Implement Benchmark Tests

**Goal:** Measure and optimize performance.

**Requirements:**
- Benchmark critical paths (JSON encoding, DB queries)
- Memory allocation benchmarks
- Compare implementations
- Profile and identify bottlenecks

**Implementation:**
```go
func BenchmarkJSONMarshal(b *testing.B) {
    task := &Task{
        ID:          1,
        Title:       "Benchmark Task",
        Description: "This is a benchmark",
        Status:      "pending",
        Priority:    "high",
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        json.Marshal(task)
    }
}

func BenchmarkJSONMarshalWithPool(b *testing.B) {
    task := &Task{
        ID:          1,
        Title:       "Benchmark Task",
        Description: "This is a benchmark",
        Status:      "pending",
        Priority:    "high",
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        buf := bufferPool.Get().(*bytes.Buffer)
        json.NewEncoder(buf).Encode(task)
        buf.Reset()
        bufferPool.Put(buf)
    }
}

// Memory allocation benchmark
func BenchmarkCreateTask(b *testing.B) {
    b.ReportAllocs() // Report memory allocations

    for i := 0; i < b.N; i++ {
        task := &Task{
            Title:    "Test",
            Priority: "high",
        }
        _ = task
    }
}

// Run: go test -bench=. -benchmem
// Output:
// BenchmarkJSONMarshal-8          1000000    1050 ns/op    256 B/op    4 allocs/op
// BenchmarkJSONMarshalWithPool-8  2000000     800 ns/op     64 B/op    1 allocs/op
```

**Concepts to Learn:**
- Writing benchmarks
- `b.N` and benchmark iteration
- `b.ReportAllocs()` for memory
- `b.ResetTimer()` for setup exclusion
- Interpreting benchmark results

---

## Advanced Patterns

### A1. Implement Circuit Breaker Pattern

**Goal:** Prevent cascading failures when external services are down.

**Requirements:**
- States: Closed (normal), Open (failing), Half-Open (testing)
- Track failure count and success count
- Configurable thresholds
- Automatic recovery attempts
- Metrics for monitoring

**Implementation:**
```go
type State int

const (
    StateClosed State = iota
    StateOpen
    StateHalfOpen
)

type CircuitBreaker struct {
    mu              sync.RWMutex
    state           State
    failures        int
    successes       int
    lastFailure     time.Time
    failureThreshold int
    successThreshold int
    timeout         time.Duration
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
    if !cb.canExecute() {
        return ErrCircuitOpen
    }

    err := fn()
    cb.recordResult(err)
    return err
}

func (cb *CircuitBreaker) canExecute() bool {
    cb.mu.RLock()
    defer cb.mu.RUnlock()

    switch cb.state {
    case StateClosed:
        return true
    case StateOpen:
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.mu.RUnlock()
            cb.mu.Lock()
            cb.state = StateHalfOpen
            cb.mu.Unlock()
            cb.mu.RLock()
            return true
        }
        return false
    case StateHalfOpen:
        return true
    }
    return false
}

func (cb *CircuitBreaker) recordResult(err error) {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()
        cb.successes = 0

        if cb.failures >= cb.failureThreshold {
            cb.state = StateOpen
        }
    } else {
        cb.successes++

        if cb.state == StateHalfOpen && cb.successes >= cb.successThreshold {
            cb.state = StateClosed
            cb.failures = 0
        }
    }
}

// Usage
var externalAPIBreaker = NewCircuitBreaker(5, 2, 30*time.Second)

func CallExternalAPI() error {
    return externalAPIBreaker.Execute(func() error {
        // actual API call
        return http.Get("https://api.example.com/data")
    })
}
```

**Concepts to Learn:**
- State machine implementation
- Failure isolation
- Automatic recovery
- Protecting against cascading failures

---

### A2. Implement Fan-Out/Fan-In Pattern

**Goal:** Process data through multiple stages concurrently.

**Requirements:**
- Pipeline with multiple stages
- Fan-out: distribute work to multiple workers
- Fan-in: collect results from multiple workers
- Proper cancellation through entire pipeline
- No goroutine leaks

**Implementation:**
```go
func Pipeline(ctx context.Context, tasks []Task) <-chan Result {
    // Stage 1: Generate
    tasksChan := generate(ctx, tasks)

    // Stage 2: Fan-out to workers
    numWorkers := 3
    workerChans := make([]<-chan Result, numWorkers)
    for i := 0; i < numWorkers; i++ {
        workerChans[i] = process(ctx, tasksChan)
    }

    // Stage 3: Fan-in results
    return merge(ctx, workerChans...)
}

func generate(ctx context.Context, tasks []Task) <-chan Task {
    out := make(chan Task)
    go func() {
        defer close(out)
        for _, task := range tasks {
            select {
            case out <- task:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

func process(ctx context.Context, in <-chan Task) <-chan Result {
    out := make(chan Result)
    go func() {
        defer close(out)
        for task := range in {
            select {
            case <-ctx.Done():
                return
            default:
                result := doWork(task)
                select {
                case out <- result:
                case <-ctx.Done():
                    return
                }
            }
        }
    }()
    return out
}

func merge(ctx context.Context, channels ...<-chan Result) <-chan Result {
    var wg sync.WaitGroup
    out := make(chan Result)

    output := func(ch <-chan Result) {
        defer wg.Done()
        for result := range ch {
            select {
            case out <- result:
            case <-ctx.Done():
                return
            }
        }
    }

    wg.Add(len(channels))
    for _, ch := range channels {
        go output(ch)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

**Concepts to Learn:**
- Pipeline pattern
- Generator pattern (function returning channel)
- Fan-out/fan-in
- Proper channel closing
- Cancellation propagation

---

### A3. Implement Retry with Exponential Backoff

**Goal:** Gracefully handle transient failures.

**Requirements:**
- Configurable max retries
- Exponential backoff with jitter
- Context support for cancellation
- Retry only on specific errors
- Logging of retry attempts

**Implementation:**
```go
type RetryConfig struct {
    MaxRetries  int
    InitialWait time.Duration
    MaxWait     time.Duration
    Multiplier  float64
}

func Retry(ctx context.Context, cfg RetryConfig, fn func() error) error {
    var lastErr error
    wait := cfg.InitialWait

    for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
        if attempt > 0 {
            // Add jitter: ±25%
            jitter := wait / 4
            actualWait := wait + time.Duration(rand.Int63n(int64(jitter*2))) - jitter

            select {
            case <-time.After(actualWait):
            case <-ctx.Done():
                return ctx.Err()
            }

            // Exponential backoff
            wait = time.Duration(float64(wait) * cfg.Multiplier)
            if wait > cfg.MaxWait {
                wait = cfg.MaxWait
            }
        }

        err := fn()
        if err == nil {
            return nil
        }

        if !isRetryable(err) {
            return err
        }

        lastErr = err
        log.Printf("attempt %d failed: %v, retrying in %v", attempt+1, err, wait)
    }

    return fmt.Errorf("max retries exceeded: %w", lastErr)
}

func isRetryable(err error) bool {
    // Network errors, timeouts, 5xx responses are retryable
    // 4xx errors are not retryable
    var netErr net.Error
    if errors.As(err, &netErr) && netErr.Timeout() {
        return true
    }
    // Add more checks...
    return true
}

// Usage
err := Retry(ctx, RetryConfig{
    MaxRetries:  3,
    InitialWait: 100 * time.Millisecond,
    MaxWait:     5 * time.Second,
    Multiplier:  2.0,
}, func() error {
    return callExternalAPI()
})
```

**Concepts to Learn:**
- Exponential backoff algorithm
- Jitter for thundering herd prevention
- Context-aware waiting
- Error classification for retry decisions

---

## Bonus Challenges

### 12. Add Task Export (CSV/JSON)

Export tasks to downloadable files.
- `GET /api/v1/tasks/export?format=csv`
- `GET /api/v1/tasks/export?format=json`

### 13. Add Bulk Operations

Perform actions on multiple tasks at once.
- `POST /api/v1/tasks/bulk/delete` - delete multiple tasks
- `POST /api/v1/tasks/bulk/status` - update status of multiple tasks

### 14. Add Task Due Date Reminders

Background job that checks for upcoming due dates and sends notifications.

### 15. Add API Rate Limiting Per User

Different rate limits based on user role (admin gets higher limits).

### 16. Add Request/Response Logging to Database

Store API request logs for debugging and analytics.

---

## Tips for Success

1. **Start with tests** - Write tests before or alongside your implementation
2. **Follow existing patterns** - Look at how similar features are implemented
3. **One feature at a time** - Complete one task before starting another
4. **Use the docs** - Reference the learning guides in `/docs/`
5. **Check the Makefile** - Use `make test`, `make lint` frequently
6. **Commit often** - Make small, focused commits

---

## Suggested Learning Paths

### Path A: Feature Development (API/CRUD Focus)
1. **Task 2 (Password Change)** - simple request handling
2. **Task 3 (Profile Update)** - error handling, constraints
3. **Task 1 (Tags)** - many-to-many relationships
4. **Task 5 (Subtasks)** - one-to-many, computed fields
5. **Task 7 (Soft Delete Recovery)** - GORM advanced features
6. **Task 4 (Comments)** - nested resources
7. **Task 6 (Activity Log)** - hooks and events
8. Advanced tasks (8-11) based on interest

### Path B: Go Core Concepts (Senior Interview Focus)
1. **G1 (Worker Pool)** - fundamental goroutine + channel pattern
2. **C2 (Rate Limiter)** - select, ticker, channel as semaphore
3. **S1 (Thread-Safe Cache)** - sync.RWMutex deep understanding
4. **CTX1 (Request Tracing)** - context propagation
5. **I1 (Repository Interfaces)** - interface design
6. **E1 (Error Handling)** - production error patterns
7. **T1-T3 (Testing)** - table tests, integration, benchmarks
8. **C1 (Event Stream)** - pub/sub with channels
9. **C3 (Singleflight)** - advanced sync patterns
10. **GEN1-GEN2 (Generics)** - modern Go
11. **A1-A3 (Advanced)** - circuit breaker, fan-out/fan-in, retry

### Path C: Full Stack (Combine Both)
Do Path A tasks 1-3, then alternate with Path B concepts.

---

## Interview Topics Covered

| Topic | Tasks | Key Concepts |
|-------|-------|--------------|
| **Goroutines** | G1, G2, C1, A2 | Lifecycle, coordination, leaks |
| **Channels** | G1, C1, C2, C3, A2 | Buffered/unbuffered, directions, closing |
| **Select** | C1, C2, CTX2 | Multiplexing, non-blocking, timeouts |
| **Context** | CTX1, CTX2, G2 | Cancellation, deadlines, values |
| **Sync Package** | S1, S2, S3, C3 | Mutex, RWMutex, Pool, Once, WaitGroup |
| **Interfaces** | I1, I2, I3 | Design, composition, mocking |
| **Error Handling** | E1, E2 | Wrapping, Is/As, panics |
| **Generics** | GEN1, GEN2 | Type parameters, constraints |
| **Testing** | T1, T2, T3 | Table-driven, integration, benchmarks |
| **Patterns** | A1, A2, A3 | Circuit breaker, pipeline, retry |

---

## Common Interview Questions These Tasks Prepare You For

### Concurrency
- "How would you process 1 million records concurrently?" → G1, A2
- "Explain goroutines vs threads" → G1, G2
- "What's the difference between buffered and unbuffered channels?" → C1, C2
- "How do you prevent goroutine leaks?" → CTX2, A2
- "Explain the select statement" → C1, C2

### Synchronization
- "When would you use Mutex vs RWMutex?" → S1
- "How does sync.Once work?" → S3
- "Explain sync.Pool and when to use it" → S2
- "How do you handle race conditions?" → S1, G1

### Error Handling
- "How do you implement error wrapping?" → E1
- "Explain errors.Is vs errors.As" → E1
- "How do you handle panics in production?" → E2

### Design Patterns
- "Implement a worker pool" → G1
- "Explain the circuit breaker pattern" → A1
- "How would you implement rate limiting?" → C2
- "Design a pub/sub system" → C1

### Testing
- "How do you write table-driven tests?" → T1
- "How do you test concurrent code?" → T2, G1
- "How do you benchmark Go code?" → T3

### System Design
- "Design a task queue system" → G1, A2
- "How would you handle retries with backoff?" → A3
- "Implement request deduplication" → C3

---

## Quick Reference: Task Categories

```
FEATURE TASKS (1-11)           GO CONCEPTS (G/C/CTX/I/E/S/GEN/T/A)
├── Beginner                   ├── Goroutines (G1, G2)
│   ├── 1. Tags               ├── Channels (C1, C2, C3)
│   ├── 2. Password Change    ├── Context (CTX1, CTX2)
│   └── 3. Profile Update     ├── Interfaces (I1, I2, I3)
├── Intermediate              ├── Errors (E1, E2)
│   ├── 4. Comments           ├── Sync (S1, S2, S3)
│   ├── 5. Subtasks           ├── Generics (GEN1, GEN2)
│   ├── 6. Activity Log       ├── Testing (T1, T2, T3)
│   └── 7. Soft Delete        └── Advanced (A1, A2, A3)
└── Advanced
    ├── 8. Redis Cache
    ├── 9. WebSocket
    ├── 10. Email Verify
    └── 11. RBAC
```

---

Good luck! Each task builds real skills for production Go development and senior-level interviews.
