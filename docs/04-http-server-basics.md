# 04 - HTTP Server Basics

## Starting a Server

### Node.js/Express vs Go

**Express:**
```javascript
const express = require('express');
const app = express();

app.get('/hello', (req, res) => {
  res.json({ message: 'Hello World' });
});

app.listen(3000, () => {
  console.log('Server running on port 3000');
});
```

**Go (standard library):**
```go
package main

import (
    "fmt"
    "net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `{"message": "Hello World"}`)
}

func main() {
    http.HandleFunc("/hello", helloHandler)

    fmt.Println("Server running on port 8080")
    http.ListenAndServe(":8080", nil)
}
```

---

## HTTP Handlers

### Handler Function Signature

Every HTTP handler in Go has this signature:

```go
func(w http.ResponseWriter, r *http.Request)
```

| Parameter | Type | Purpose | Express Equivalent |
|-----------|------|---------|-------------------|
| `w` | `http.ResponseWriter` | Write response | `res` |
| `r` | `*http.Request` | Read request | `req` |

### Writing Responses

**Express:**
```javascript
// Text
res.send('Hello');

// JSON
res.json({ message: 'Hello' });

// Status + JSON
res.status(201).json({ id: 1 });

// Headers
res.set('X-Custom-Header', 'value');
```

**Go:**
```go
// Text
fmt.Fprintf(w, "Hello")
// or
w.Write([]byte("Hello"))

// JSON
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]string{"message": "Hello"})

// Status + JSON
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated)  // Must be before Write!
json.NewEncoder(w).Encode(map[string]int{"id": 1})

// Headers
w.Header().Set("X-Custom-Header", "value")
```

**Important:** In Go, set headers BEFORE calling `WriteHeader()` or `Write()`.

### Reading Requests

**Express:**
```javascript
// URL parameters
const id = req.params.id;

// Query string
const page = req.query.page;

// Body (with body-parser)
const { name } = req.body;

// Headers
const token = req.headers.authorization;
```

**Go:**
```go
// URL path (need a router like chi/gorilla for path params)
// Standard library doesn't support /users/:id patterns

// Query string
page := r.URL.Query().Get("page")

// Body
var body struct {
    Name string `json:"name"`
}
json.NewDecoder(r.Body).Decode(&body)
defer r.Body.Close()  // Always close the body!

// Headers
token := r.Header.Get("Authorization")
```

---

## JSON Handling

### Encoding (Struct → JSON)

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age,omitempty"` // omitempty: skip if zero value
}

func handler(w http.ResponseWriter, r *http.Request) {
    user := User{ID: 1, Name: "John"}

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
    // Output: {"id":1,"name":"John"}
    // Note: age is omitted because it's 0 and has omitempty
}
```

### Decoding (JSON → Struct)

```go
func createUser(w http.ResponseWriter, r *http.Request) {
    var user User

    // Decode request body
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    fmt.Printf("Created user: %+v\n", user)
}
```

### JSON Struct Tags

```go
type User struct {
    ID        int    `json:"id"`                    // Rename to "id"
    FirstName string `json:"first_name"`            // Rename to "first_name"
    Password  string `json:"-"`                     // Never include in JSON
    Age       int    `json:"age,omitempty"`         // Omit if zero
    Email     string `json:"email,omitempty,string"` // Multiple options
}
```

---

## Routing

### Standard Library (Basic)

```go
// Exact path matching
http.HandleFunc("/hello", helloHandler)
http.HandleFunc("/users", usersHandler)

// Note: "/" matches everything!
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    fmt.Fprintf(w, "Home")
})
```

**Limitations:**
- No path parameters (`/users/:id`)
- No method-specific routing (GET vs POST)
- No middleware chain

### Popular Routers

For production apps, use a router like:

**chi** (lightweight, idiomatic):
```go
import "github.com/go-chi/chi/v5"

r := chi.NewRouter()
r.Get("/users/{id}", getUserHandler)
r.Post("/users", createUserHandler)
```

**gorilla/mux** (feature-rich):
```go
import "github.com/gorilla/mux"

r := mux.NewRouter()
r.HandleFunc("/users/{id}", getUserHandler).Methods("GET")
```

**gin** (fast, feature-rich):
```go
import "github.com/gin-gonic/gin"

r := gin.Default()
r.GET("/users/:id", getUser)
```

---

## Middleware

### Concept

Middleware wraps handlers to add functionality (logging, auth, CORS).

**Express:**
```javascript
app.use(express.json());       // Parse JSON bodies
app.use(cors());               // Enable CORS
app.use(authMiddleware);       // Custom auth

function authMiddleware(req, res, next) {
    if (req.headers.authorization) {
        next();
    } else {
        res.status(401).send('Unauthorized');
    }
}
```

**Go (standard library pattern):**
```go
// Middleware is a function that wraps a handler
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next(w, r)  // Call the next handler
    }
}

// Usage
http.HandleFunc("/hello", loggingMiddleware(helloHandler))
```

**Go (with chi):**
```go
r := chi.NewRouter()
r.Use(middleware.Logger)      // Built-in logger
r.Use(middleware.Recoverer)   // Panic recovery
r.Use(corsMiddleware)         // Custom middleware

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        next.ServeHTTP(w, r)
    })
}
```

---

## Error Handling

### Express vs Go

**Express:**
```javascript
app.get('/users/:id', async (req, res, next) => {
    try {
        const user = await findUser(req.params.id);
        if (!user) {
            return res.status(404).json({ error: 'User not found' });
        }
        res.json(user);
    } catch (err) {
        next(err);  // Pass to error handler middleware
    }
});

// Global error handler
app.use((err, req, res, next) => {
    res.status(500).json({ error: 'Internal server error' });
});
```

**Go:**
```go
func getUserHandler(w http.ResponseWriter, r *http.Request) {
    user, err := findUser(id)
    if err != nil {
        // Handle specific error types
        if errors.Is(err, ErrNotFound) {
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }
        // Generic error
        log.Printf("Error: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(user)
}
```

**Key Difference:** Go doesn't have a global error handler. Each handler manages its own errors explicitly.

---

## Full Example: RESTful Endpoint

Here's a complete example showing the patterns together:

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
)

type Task struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Status string `json:"status"`
}

var tasks = []Task{
    {ID: 1, Title: "Learn Go", Status: "in_progress"},
    {ID: 2, Title: "Build API", Status: "pending"},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tasks)
}

func main() {
    http.HandleFunc("/tasks", getTasks)

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## Exercises

1. Add a new endpoint `/tasks` that returns a list of tasks
2. Add query parameter support (e.g., `/tasks?status=pending`)
3. Try adding a simple logging middleware
4. Implement POST `/tasks` to create a new task (hint: decode JSON body)

---

## Next Steps

- Read **05-next-steps.md** to see what we'll learn next
- Try implementing CRUD operations for tasks
- Explore chi or gorilla/mux for better routing
