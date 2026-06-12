# 13 - Middleware Patterns

Middleware are functions that run before/after handlers. They're perfect for cross-cutting concerns like logging, authentication, and error handling.

## Middleware Concept

```
Request → [Middleware1] → [Middleware2] → [Handler] → Response
              ↓              ↓               ↓
           Before         Before          Execute
              ↓              ↓               ↓
            Next()         Next()         Return
              ↓              ↓               ↓
            After          After           Done
```

---

## Express vs Gin Middleware

### Express.js

```javascript
// Middleware function
const logger = (req, res, next) => {
  console.log(`${req.method} ${req.path}`);
  next();  // Continue to next middleware/handler
};

// Apply globally
app.use(logger);

// Apply to specific route
app.get('/users', authMiddleware, getUsers);
```

### Gin

```go
// Middleware function
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Printf("%s %s\n", c.Request.Method, c.Request.URL.Path)
        c.Next()  // Continue to next middleware/handler
    }
}

// Apply globally
r.Use(Logger())

// Apply to specific route
r.GET("/users", AuthMiddleware(), getUsers)

// Apply to route group
api := r.Group("/api")
api.Use(AuthMiddleware())
```

---

## Our Middleware Stack

See `cmd/api/main.go`:

```go
r := gin.New()

// Global middleware (order matters!)
r.Use(middleware.Recovery())    // 1. Catch panics
r.Use(middleware.Logger())      // 2. Log requests
r.Use(middleware.CORS())        // 3. Handle CORS

// Route-specific middleware
protected := r.Group("/api/v1")
protected.Use(middleware.Auth(authService))  // Auth required
{
    protected.GET("/tasks", taskHandler.List)
}
```

---

## 1. Logger Middleware

See `internal/middleware/logger.go`:

```go
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Before handler
        start := time.Now()
        path := c.Request.URL.Path
        method := c.Request.Method

        // Process request
        c.Next()

        // After handler
        latency := time.Since(start)
        statusCode := c.Writer.Status()

        logger.Info("Request",
            zap.String("method", method),
            zap.String("path", path),
            zap.Int("status", statusCode),
            zap.Duration("latency", latency),
            zap.String("ip", c.ClientIP()),
        )
    }
}
```

**Key Points:**
- `c.Next()` calls the next handler in chain
- Code before `Next()` runs before handler
- Code after `Next()` runs after handler
- Can access response info (status, size) after `Next()`

---

## 2. Recovery Middleware

See `internal/middleware/recovery.go`:

```go
func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // Log the panic with stack trace
                logger.Error("Panic recovered",
                    zap.Any("error", err),
                    zap.String("stack", string(debug.Stack())),
                )

                // Return 500 error
                c.AbortWithStatusJSON(http.StatusInternalServerError,
                    dto.APIResponse[any]{
                        Success: false,
                        Error: &dto.ErrorDetail{
                            Code:    "INTERNAL_ERROR",
                            Message: "Internal server error",
                        },
                    },
                )
            }
        }()

        c.Next()
    }
}
```

**Key Points:**
- Uses `defer` to catch panics
- `recover()` catches panic values
- Must use `c.Abort()` to stop chain on error

---

## 3. CORS Middleware

See `internal/middleware/cors.go`:

```go
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")

        // Allow all origins in development
        // In production, check against allowed list
        c.Header("Access-Control-Allow-Origin", origin)
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Max-Age", "86400")

        // Handle preflight
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}
```

**Key Points:**
- Sets CORS headers on every response
- Handles OPTIONS preflight requests
- `c.AbortWithStatus()` stops chain without error response

---

## 4. Auth Middleware

See `internal/middleware/auth.go`:

```go
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            utils.UnauthorizedResponse(c, "Authorization header required")
            c.Abort()
            return
        }

        // Parse "Bearer <token>"
        if !strings.HasPrefix(authHeader, "Bearer ") {
            utils.UnauthorizedResponse(c, "Invalid authorization format")
            c.Abort()
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")

        // Validate token
        claims, err := authService.ValidateAccessToken(token)
        if err != nil {
            utils.UnauthorizedResponse(c, "Invalid or expired token")
            c.Abort()
            return
        }

        // Set user info in context for handlers
        c.Set("userID", claims.UserID)
        c.Set("user", claims)

        c.Next()
    }
}
```

**Key Points:**
- Uses `c.Set()` to pass data to handlers
- Uses `c.Abort()` to stop chain on auth failure
- Returns appropriate HTTP response before aborting

---

## 5. Rate Limit Middleware

See `internal/middleware/ratelimit.go`:

```go
type RateLimiter struct {
    requests map[string]*clientInfo
    mu       sync.RWMutex  // Thread-safe access
    limit    int           // Max requests
    window   time.Duration // Time window
}

type clientInfo struct {
    count     int
    expiresAt time.Time
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
    rl := &RateLimiter{
        requests: make(map[string]*clientInfo),
        limit:    limit,
        window:   window,
    }

    // Start cleanup goroutine
    go rl.cleanupLoop()

    return rl
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        clientIP := c.ClientIP()

        if !rl.Allow(clientIP) {
            c.Header("Retry-After", fmt.Sprintf("%d", int(rl.window.Seconds())))
            utils.TooManyRequestsResponse(c, "Rate limit exceeded")
            c.Abort()
            return
        }

        c.Next()
    }
}

func (rl *RateLimiter) Allow(clientIP string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    now := time.Now()
    info, exists := rl.requests[clientIP]

    if !exists || now.After(info.expiresAt) {
        // New window
        rl.requests[clientIP] = &clientInfo{
            count:     1,
            expiresAt: now.Add(rl.window),
        }
        return true
    }

    if info.count >= rl.limit {
        return false  // Rate limited
    }

    info.count++
    return true
}
```

**Key Points:**
- Uses `sync.RWMutex` for thread-safe map access
- Background goroutine cleans up expired entries
- Sets `Retry-After` header for rate-limited responses

---

## Middleware Order Matters!

```go
// Order: Recovery → Logger → CORS → Auth → Handler

// If you put Logger before Recovery:
// - Panics won't be logged properly

// If you put Auth before CORS:
// - Preflight requests will fail auth

// Correct order:
r.Use(middleware.Recovery())  // Always first
r.Use(middleware.Logger())    // Log everything
r.Use(middleware.CORS())      // Before auth (for preflight)
// Auth is applied to specific routes, not globally
```

---

## Creating Custom Middleware

### Template

```go
func MyMiddleware(/* dependencies */) gin.HandlerFunc {
    // Setup code runs once when middleware is registered

    return func(c *gin.Context) {
        // Before handler
        // - Read request
        // - Validate
        // - Set context values

        c.Next()  // Call next handler

        // After handler
        // - Read response
        // - Modify headers
        // - Log
    }
}
```

### Example: Request ID Middleware

```go
func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Check if request ID was provided
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }

        // Set in context for logging
        c.Set("requestID", requestID)

        // Set in response header
        c.Header("X-Request-ID", requestID)

        c.Next()
    }
}
```

### Example: Timeout Middleware

```go
func Timeout(timeout time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Create context with timeout
        ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
        defer cancel()

        // Replace request context
        c.Request = c.Request.WithContext(ctx)

        // Channel to track completion
        done := make(chan struct{})

        go func() {
            c.Next()
            close(done)
        }()

        select {
        case <-done:
            // Handler completed
        case <-ctx.Done():
            // Timeout
            c.AbortWithStatusJSON(http.StatusGatewayTimeout,
                gin.H{"error": "Request timeout"})
        }
    }
}
```

---

## Middleware Patterns

### 1. Skip for Certain Paths

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Skip auth for public paths
        publicPaths := []string{"/health", "/ready", "/live"}
        for _, path := range publicPaths {
            if c.Request.URL.Path == path {
                c.Next()
                return
            }
        }

        // Normal auth logic
        // ...
    }
}
```

### 2. Conditional Middleware

```go
func ConditionalMiddleware(condition bool, mw gin.HandlerFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        if condition {
            mw(c)
        } else {
            c.Next()
        }
    }
}

// Usage
if config.EnableRateLimit {
    r.Use(rateLimiter.Middleware())
}
```

### 3. Middleware Composition

```go
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        c.Next()
    }
}

// Combine multiple security middlewares
func Security() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Apply all security measures
        SecurityHeaders()(c)
        // Could add more...
    }
}
```

---

## Comparison: Express vs Gin Middleware

| Aspect | Express | Gin |
|--------|---------|-----|
| Signature | `(req, res, next)` | `func(c *gin.Context)` |
| Continue | `next()` | `c.Next()` |
| Stop | Don't call `next()` | `c.Abort()` |
| Stop with error | `next(err)` | `c.AbortWithStatus()` |
| Set value | `req.user = user` | `c.Set("user", user)` |
| Get value | `req.user` | `c.Get("user")` |
| Error handler | Separate `(err, req, res, next)` | Recovery middleware |

---

## Exercises

1. Look at `internal/middleware/logger.go` - what information is logged?
2. Look at `internal/middleware/ratelimit.go` - how does the cleanup work?
3. Create a middleware that adds a `X-Response-Time` header
4. Modify the auth middleware to support API keys in addition to JWT
