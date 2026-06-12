# 10 - Goroutines and Channels

This is where Go shines! Goroutines and channels are Go's powerful concurrency primitives.

## Goroutines

### What is a Goroutine?

A goroutine is a lightweight thread managed by the Go runtime. They're incredibly cheap - you can run thousands without issues.

**Node.js (async/await):**
```javascript
async function sendEmail(email) {
  await emailService.send(email);
}

// This is still single-threaded!
await sendEmail('user@example.com');
```

**Go (goroutine):**
```go
func sendEmail(email string) {
    // email sending logic
}

// Run in a separate goroutine (truly concurrent!)
go sendEmail("user@example.com")
```

### Key Differences from Node.js

| Aspect | Node.js | Go |
|--------|---------|-----|
| Threading | Single-threaded event loop | Multiple OS threads |
| Concurrency | Async I/O | Goroutines (M:N scheduling) |
| CPU-bound tasks | Blocks event loop | Runs on separate core |
| Memory per "thread" | N/A | ~2KB per goroutine |
| Max concurrent | Limited by event loop | Millions of goroutines |

### Creating Goroutines

```go
// Simple goroutine
go myFunction()

// Anonymous function
go func() {
    fmt.Println("Running in goroutine")
}()

// With parameters (careful with closures!)
for i := 0; i < 10; i++ {
    go func(n int) {  // Pass i as parameter!
        fmt.Println(n)
    }(i)
}
```

### Real Example: Our NotificationService

See `internal/services/notification_service.go`:

```go
// Worker pool pattern
func (s *NotificationService) startWorkers() {
    for i := 0; i < s.workerCount; i++ {
        s.wg.Add(1)
        go func(workerID int) {
            defer s.wg.Done()
            s.workerLoop(workerID)
        }(i)
    }
}
```

---

## Channels

### What is a Channel?

Channels are pipes for communication between goroutines. They're typed and thread-safe.

**Node.js (EventEmitter):**
```javascript
const emitter = new EventEmitter();

// Producer
emitter.emit('message', 'Hello');

// Consumer
emitter.on('message', (msg) => console.log(msg));
```

**Go (Channel):**
```go
// Create channel
messages := make(chan string)

// Producer (in goroutine)
go func() {
    messages <- "Hello"  // Send to channel
}()

// Consumer
msg := <-messages  // Receive from channel (blocks!)
fmt.Println(msg)
```

### Buffered vs Unbuffered Channels

```go
// Unbuffered - sender blocks until receiver is ready
ch := make(chan int)

// Buffered - sender can queue up to N items
ch := make(chan int, 10)
```

| Type | Behavior | Use Case |
|------|----------|----------|
| Unbuffered | Synchronous, blocks | Handoff, synchronization |
| Buffered | Async up to capacity | Queues, decoupling |

### Channel Operations

```go
ch := make(chan int, 5)

// Send
ch <- 42

// Receive
value := <-ch

// Receive with ok (detect closed channel)
value, ok := <-ch
if !ok {
    // Channel was closed
}

// Close channel
close(ch)

// Range over channel (until closed)
for value := range ch {
    fmt.Println(value)
}
```

### Select Statement

`select` is like switch but for channels:

```go
select {
case msg := <-messages:
    fmt.Println("Received:", msg)
case <-time.After(5 * time.Second):
    fmt.Println("Timeout!")
case <-ctx.Done():
    fmt.Println("Cancelled")
default:
    fmt.Println("No message ready")
}
```

---

## Common Patterns

### 1. Worker Pool

Multiple workers consuming from a shared channel:

```go
func workerPool(jobs <-chan Job, results chan<- Result, workers int) {
    var wg sync.WaitGroup

    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for job := range jobs {
                results <- process(job)
            }
        }(i)
    }

    wg.Wait()
    close(results)
}
```

### 2. Fan-Out/Fan-In

Distribute work and collect results:

```go
// Fan-out: Multiple goroutines reading from one channel
func fanOut(input <-chan int, workers int) []<-chan int {
    outputs := make([]<-chan int, workers)
    for i := 0; i < workers; i++ {
        outputs[i] = worker(input)
    }
    return outputs
}

// Fan-in: Merge multiple channels into one
func fanIn(inputs ...<-chan int) <-chan int {
    output := make(chan int)
    var wg sync.WaitGroup

    for _, ch := range inputs {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for v := range c {
                output <- v
            }
        }(ch)
    }

    go func() {
        wg.Wait()
        close(output)
    }()

    return output
}
```

### 3. Timeout Pattern

```go
func fetchWithTimeout(url string, timeout time.Duration) ([]byte, error) {
    result := make(chan []byte, 1)
    errCh := make(chan error, 1)

    go func() {
        data, err := fetch(url)
        if err != nil {
            errCh <- err
            return
        }
        result <- data
    }()

    select {
    case data := <-result:
        return data, nil
    case err := <-errCh:
        return nil, err
    case <-time.After(timeout):
        return nil, errors.New("timeout")
    }
}
```

---

## WaitGroup

Wait for multiple goroutines to complete:

```go
var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)  // Increment counter
    go func(n int) {
        defer wg.Done()  // Decrement when done
        doWork(n)
    }(i)
}

wg.Wait()  // Block until counter is 0
```

**Node.js equivalent:**
```javascript
await Promise.all([
    doWork(0),
    doWork(1),
    doWork(2),
    doWork(3),
    doWork(4),
]);
```

---

## Context for Cancellation

Context propagates cancellation through call chains:

```go
// Create context with cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Create context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Check for cancellation
select {
case <-ctx.Done():
    return ctx.Err()
default:
    // Continue work
}
```

---

## Best Practices

### DO:
- Always close channels from the sender side
- Use `defer wg.Done()` to ensure cleanup
- Pass context to long-running operations
- Use buffered channels when producer/consumer speeds differ

### DON'T:
- Close a channel from the receiver side
- Send on a closed channel (panics!)
- Forget to handle the `ok` value from channel receive
- Create goroutines without a way to stop them

---

## Exercises

1. Look at `notification_service.go` and trace the flow of a notification
2. Add a method to drain all pending notifications on shutdown
3. Create a simple worker pool that processes numbers
4. Add timeout to a database query using context
