# 05 - Next Steps

Now that you have a basic Go HTTP server running, here's what we'll explore next. These are the features that make Go unique and powerful!

---

## Coming Up: Loops

### Basic Loops

Go has only ONE loop keyword: `for`. It does everything!

```go
// Classic for loop (like JavaScript for)
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// While loop (for without init/post)
count := 0
for count < 10 {
    fmt.Println(count)
    count++
}

// Infinite loop
for {
    fmt.Println("Forever!")
    break  // Don't actually run forever
}

// Range loop (like for...of in JavaScript)
tasks := []string{"Learn Go", "Build API", "Deploy"}
for index, task := range tasks {
    fmt.Printf("%d: %s\n", index, task)
}

// Range with maps
ages := map[string]int{"Alice": 30, "Bob": 25}
for name, age := range ages {
    fmt.Printf("%s is %d\n", name, age)
}
```

**JavaScript comparison:**
```javascript
// for
for (let i = 0; i < 10; i++) { }

// while
while (count < 10) { }

// for...of
for (const task of tasks) { }

// for...in (objects)
for (const key in object) { }
```

---

## Coming Up: Goroutines

Goroutines are Go's killer feature for concurrency. They're like lightweight threads.

### Basic Goroutine

```go
// Normal function call (blocks)
doWork()  // Waits for completion

// Goroutine (runs concurrently)
go doWork()  // Returns immediately, work happens in background
```

**JavaScript comparison:**
```javascript
// Promises/async (still single-threaded)
doWork();  // Sync
await doWork();  // Async, but one at a time

// Go goroutines run truly in parallel on multiple CPU cores!
```

### Real Example: Concurrent HTTP Requests

```go
func fetchAll(urls []string) {
    for _, url := range urls {
        go fetch(url)  // Each runs concurrently
    }
}

func fetch(url string) {
    resp, _ := http.Get(url)
    fmt.Printf("Fetched %s: %d\n", url, resp.StatusCode)
}
```

---

## Coming Up: Channels

Channels are how goroutines communicate safely. Think of them as pipes for data.

### Basic Channel

```go
// Create a channel
messages := make(chan string)

// Send to channel (in a goroutine)
go func() {
    messages <- "Hello!"  // Send
}()

// Receive from channel (blocks until data arrives)
msg := <-messages  // Receive
fmt.Println(msg)   // "Hello!"
```

### Real Example: Collecting Results

```go
func fetchAllWithResults(urls []string) []string {
    results := make(chan string, len(urls))  // Buffered channel

    for _, url := range urls {
        go func(u string) {
            resp, _ := http.Get(u)
            results <- fmt.Sprintf("%s: %d", u, resp.StatusCode)
        }(url)
    }

    // Collect all results
    var allResults []string
    for i := 0; i < len(urls); i++ {
        allResults = append(allResults, <-results)
    }
    return allResults
}
```

**JavaScript comparison:**
```javascript
// Similar to Promise.all()
const results = await Promise.all(urls.map(url => fetch(url)));
```

---

## Coming Up: Structs & Methods

We've used structs, but there's more to learn!

### Methods on Structs

```go
type Task struct {
    ID     int
    Title  string
    Status string
}

// Method with value receiver
func (t Task) IsComplete() bool {
    return t.Status == "complete"
}

// Method with pointer receiver (can modify the struct)
func (t *Task) Complete() {
    t.Status = "complete"
}

// Usage
task := Task{ID: 1, Title: "Learn Go", Status: "pending"}
task.Complete()              // Modifies task
fmt.Println(task.IsComplete())  // true
```

**JavaScript comparison:**
```javascript
class Task {
    constructor(id, title, status) {
        this.id = id;
        this.title = title;
        this.status = status;
    }

    isComplete() {
        return this.status === 'complete';
    }

    complete() {
        this.status = 'complete';
    }
}
```

---

## Coming Up: Interfaces

Go interfaces are implicit - no `implements` keyword!

```go
// Define an interface
type TaskStore interface {
    GetAll() []Task
    GetByID(id int) (*Task, error)
    Create(task Task) error
}

// Any type with these methods automatically implements it
type MemoryStore struct {
    tasks []Task
}

func (m *MemoryStore) GetAll() []Task { ... }
func (m *MemoryStore) GetByID(id int) (*Task, error) { ... }
func (m *MemoryStore) Create(task Task) error { ... }

// MemoryStore now implements TaskStore!

// Function that accepts the interface
func ListTasks(store TaskStore) {
    tasks := store.GetAll()
    // Works with ANY implementation!
}
```

---

## Coming Up: Error Handling Patterns

More sophisticated error handling:

```go
// Custom errors
var ErrNotFound = errors.New("not found")

func GetTask(id int) (*Task, error) {
    task := findTask(id)
    if task == nil {
        return nil, ErrNotFound
    }
    return task, nil
}

// Checking error types
task, err := GetTask(1)
if errors.Is(err, ErrNotFound) {
    // Handle not found
}
```

---

## Coming Up: Database (MySQL)

We'll connect to MySQL and build real CRUD:

```go
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

db, err := sql.Open("mysql", "user:pass@tcp(localhost:3306)/taskmanager")

// Query
rows, err := db.Query("SELECT id, title FROM tasks")
for rows.Next() {
    var task Task
    rows.Scan(&task.ID, &task.Title)
}

// Insert
result, err := db.Exec("INSERT INTO tasks (title) VALUES (?)", "New Task")
id, _ := result.LastInsertId()
```

---

## Learning Path Suggestions

### Week 1-2: Basics (where we are now)
- [ ] Basic HTTP endpoints
- [ ] Structs and JSON
- [ ] Environment configuration
- [ ] Error handling basics

### Week 3-4: Data & Loops
- [ ] Slices and maps
- [ ] All loop patterns
- [ ] MySQL connection
- [ ] CRUD operations

### Week 5-6: Concurrency
- [ ] Goroutines
- [ ] Channels
- [ ] sync package (Mutex, WaitGroup)
- [ ] Context for cancellation

### Week 7-8: Production Ready
- [ ] Proper error handling
- [ ] Logging (logrus/zap)
- [ ] Configuration (viper)
- [ ] Testing

---

## Useful Resources

### Official
- [Go Tour](https://go.dev/tour/) - Interactive tutorial
- [Effective Go](https://go.dev/doc/effective_go) - Best practices
- [Go by Example](https://gobyexample.com/) - Code examples

### For Node.js Developers
- [Go for Node.js Developers](https://github.com/miguelmota/golang-for-nodejs-developers)
- Compare patterns side-by-side

### Tools to Install
```bash
# Linter
brew install golangci-lint

# Live reload (like nodemon)
go install github.com/cosmtrek/air@latest

# Better REPL
go install github.com/x-motemen/gore/cmd/gore@latest
```

---

## Quick Reference: Node.js → Go

| Concept | Node.js | Go |
|---------|---------|-----|
| Package manager | npm/yarn | go modules |
| Entry point | index.js | main.go |
| Async | Promises/async-await | Goroutines/channels |
| Array | Array | Slice |
| Object | Object | Map / Struct |
| Class | class | struct + methods |
| Interface | interface | interface (implicit) |
| Try/catch | try/catch | if err != nil |
| null/undefined | null, undefined | nil |
| Threads | Worker threads | Goroutines (built-in) |

---

Ready to continue? Let's build something! Try:
1. Adding a `/tasks` endpoint that returns static tasks
2. Implementing different loop patterns
3. Creating your first goroutine

Happy coding! 🚀
