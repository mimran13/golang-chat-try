# 06 - Gin Web Framework

Gin is Go's most popular HTTP web framework. If you know Express.js, Gin will feel familiar!

## Why Gin?

| Feature | Gin | Express.js |
|---------|-----|------------|
| Performance | Extremely fast (httprouter) | Good |
| Routing | Tree-based, efficient | Linear |
| Middleware | Built-in chain | Plugin-based |
| Validation | Integrated (validator) | Separate (joi, yup) |
| JSON | Native, fast | JSON.parse/stringify |

---

## Basic Setup

**Express.js:**
```javascript
const express = require('express');
const app = express();

app.use(express.json());

app.get('/', (req, res) => {
  res.json({ message: 'Hello' });
});

app.listen(3000);
```

**Gin:**
```go
package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()  // Includes logger and recovery middleware

    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Hello"})
    })

    r.Run(":3000")
}
```

---

## Router and Routes

### Basic Routes

```go
r := gin.Default()

// HTTP methods
r.GET("/users", getUsers)
r.POST("/users", createUser)
r.PUT("/users/:id", updateUser)
r.DELETE("/users/:id", deleteUser)
r.PATCH("/users/:id", patchUser)

// Any method
r.Any("/any", handleAny)

// Multiple methods
r.Match([]string{"GET", "POST"}, "/match", handleMatch)
```

### Route Groups

**Express.js:**
```javascript
const router = express.Router();
router.get('/users', getUsers);
router.post('/users', createUser);
app.use('/api/v1', router);
```

**Gin:**
```go
// Group routes under /api/v1
v1 := r.Group("/api/v1")
{
    // /api/v1/users
    v1.GET("/users", getUsers)
    v1.POST("/users", createUser)

    // Nested group: /api/v1/admin
    admin := v1.Group("/admin")
    {
        admin.GET("/stats", getStats)
    }
}
```

### Route Parameters

```go
// Path parameter: /users/123
r.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")  // "123"
    c.JSON(200, gin.H{"id": id})
})

// Wildcard: /files/path/to/file.txt
r.GET("/files/*filepath", func(c *gin.Context) {
    filepath := c.Param("filepath")  // "/path/to/file.txt"
    c.JSON(200, gin.H{"path": filepath})
})
```

---

## Request Handling

### Query Parameters

**Express.js:**
```javascript
app.get('/search', (req, res) => {
  const { q, page = 1 } = req.query;
});
```

**Gin:**
```go
r.GET("/search", func(c *gin.Context) {
    // Get single value
    q := c.Query("q")           // "" if not present
    q := c.DefaultQuery("q", "default")

    // Get with existence check
    page, exists := c.GetQuery("page")

    // Get array: ?tags=go&tags=web
    tags := c.QueryArray("tags")

    // Get map: ?filters[status]=active
    filters := c.QueryMap("filters")
})
```

### Request Body

**Express.js:**
```javascript
app.post('/users', (req, res) => {
  const { name, email } = req.body;
});
```

**Gin:**
```go
type CreateUserRequest struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
}

r.POST("/users", func(c *gin.Context) {
    var req CreateUserRequest

    // Bind JSON body to struct
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Use req.Name, req.Email
    c.JSON(201, gin.H{"name": req.Name})
})
```

### Binding Methods

```go
// JSON body
c.ShouldBindJSON(&obj)

// Form data
c.ShouldBind(&obj)

// Query parameters
c.ShouldBindQuery(&obj)

// URI parameters
c.ShouldBindUri(&obj)

// Headers
c.ShouldBindHeader(&obj)

// Must bind (panics on error) - avoid in production
c.BindJSON(&obj)
```

### Headers

```go
r.GET("/", func(c *gin.Context) {
    // Get header
    auth := c.GetHeader("Authorization")
    contentType := c.ContentType()

    // Set response header
    c.Header("X-Custom", "value")
})
```

---

## Response Handling

### JSON Response

```go
// Using gin.H (map shorthand)
c.JSON(200, gin.H{
    "status": "success",
    "data": user,
})

// Using struct
type Response struct {
    Status string `json:"status"`
    Data   any    `json:"data"`
}
c.JSON(200, Response{Status: "success", Data: user})

// Pretty JSON (for debugging)
c.IndentedJSON(200, data)

// Secure JSON (prevents JSON hijacking)
c.SecureJSON(200, data)

// JSONP (for cross-domain)
c.JSONP(200, data)
```

### Other Response Types

```go
// String
c.String(200, "Hello %s", name)

// XML
c.XML(200, data)

// YAML
c.YAML(200, data)

// HTML
c.HTML(200, "template.html", data)

// File
c.File("/path/to/file")

// Data (raw bytes)
c.Data(200, "text/plain", []byte("raw data"))

// Redirect
c.Redirect(302, "/new-url")
```

### Status Codes

```go
import "net/http"

c.JSON(http.StatusOK, data)           // 200
c.JSON(http.StatusCreated, data)      // 201
c.JSON(http.StatusBadRequest, err)    // 400
c.JSON(http.StatusUnauthorized, err)  // 401
c.JSON(http.StatusNotFound, err)      // 404
c.JSON(http.StatusInternalServerError, err)  // 500
```

---

## Context (gin.Context)

The `*gin.Context` is the most important type in Gin. It carries:
- Request information
- Response writer
- Middleware chain control
- Request-scoped values

### Storing Values

**Express.js:**
```javascript
app.use((req, res, next) => {
  req.user = { id: 123 };
  next();
});
```

**Gin:**
```go
// Set value in middleware
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Set("userID", uint(123))
        c.Next()
    }
}

// Get value in handler
func Handler(c *gin.Context) {
    // Get returns (value, exists)
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(401, gin.H{"error": "not authenticated"})
        return
    }

    // Type assertion
    id := userID.(uint)

    // Or use typed getters
    id := c.GetUint("userID")      // Returns 0 if not found
    name := c.GetString("name")    // Returns "" if not found
}
```

### Request Context (for cancellation)

```go
func Handler(c *gin.Context) {
    // Get Go's context.Context for passing to services
    ctx := c.Request.Context()

    // Pass to database/service calls
    user, err := userRepo.FindByID(ctx, id)

    // Check if request was cancelled
    select {
    case <-ctx.Done():
        return  // Client disconnected
    default:
        // Continue processing
    }
}
```

---

## Our Application Structure

See `cmd/api/main.go`:

```go
func main() {
    // Create router
    r := gin.New()  // Empty router (no default middleware)

    // Add middleware
    r.Use(middleware.Logger())      // Custom Zap logger
    r.Use(middleware.Recovery())    // Panic recovery
    r.Use(middleware.CORS())        // CORS headers

    // Health endpoints (no auth)
    health := handlers.NewHealthHandler(db)
    r.GET("/health", health.Health)
    r.GET("/ready", health.Ready)
    r.GET("/live", health.Live)

    // API v1 routes
    v1 := r.Group("/api/v1")
    {
        // Auth routes (public)
        authHandler.RegisterRoutes(v1, authMiddleware)

        // Task routes (protected)
        taskHandler.RegisterRoutes(v1, authMiddleware)
    }

    // Start server
    r.Run(":8080")
}
```

---

## gin.Default() vs gin.New()

```go
// gin.Default() includes:
// - Logger middleware (logs requests to stdout)
// - Recovery middleware (recovers from panics)
r := gin.Default()

// gin.New() is bare:
// - No middleware
// - You add what you need
r := gin.New()
r.Use(gin.Logger())
r.Use(gin.Recovery())

// We use gin.New() because we have custom middleware
r := gin.New()
r.Use(middleware.Logger())      // Our Zap logger
r.Use(middleware.Recovery())    // Our custom recovery
```

---

## Comparison: Express vs Gin Patterns

### Express.js Controller
```javascript
class UserController {
  async create(req, res, next) {
    try {
      const user = await userService.create(req.body);
      res.status(201).json(user);
    } catch (error) {
      next(error);
    }
  }
}

router.post('/users', userController.create);
```

### Gin Handler
```go
type UserHandler struct {
    service *UserService
}

func (h *UserHandler) Create(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    user, err := h.service.Create(c.Request.Context(), &req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(201, user)
}

r.POST("/users", userHandler.Create)
```

---

## Exercises

1. Look at `internal/handlers/auth_handler.go` - how are routes registered?
2. Look at `internal/handlers/task_handler.go` - how is pagination handled?
3. Try adding a new endpoint `/api/v1/ping` that returns `{"pong": true}`
4. Add a route parameter endpoint `/api/v1/echo/:message` that echoes back the message
