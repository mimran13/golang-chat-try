# 12 - Error Handling in Go

Error handling in Go is explicit and verbose - but this is by design. Unlike exceptions in JavaScript, errors are values that must be handled.

## Go's Error Philosophy

```
"Errors are values" - Rob Pike
```

### JavaScript (Exceptions)

```javascript
try {
  const user = await userService.findById(id);
  if (!user) throw new NotFoundException('User not found');
  return user;
} catch (error) {
  if (error instanceof NotFoundException) {
    res.status(404).json({ error: error.message });
  } else {
    res.status(500).json({ error: 'Internal server error' });
  }
}
```

### Go (Error Values)

```go
user, err := userService.FindByID(ctx, id)
if err != nil {
    return nil, err  // Propagate error
}
if user == nil {
    return nil, ErrUserNotFound  // Return specific error
}
return user, nil
```

---

## The error Interface

```go
// error is a built-in interface
type error interface {
    Error() string
}

// Any type with Error() method is an error
type MyError struct {
    Message string
}

func (e *MyError) Error() string {
    return e.Message
}
```

---

## Creating Errors

### Simple Errors

```go
import "errors"

// Create simple error
err := errors.New("something went wrong")

// Create formatted error
import "fmt"
err := fmt.Errorf("user %d not found", userID)
```

### Sentinel Errors (Error Constants)

See `internal/errors/errors.go`:

```go
package apperrors

import "errors"

// Sentinel errors - define once, compare everywhere
var (
    ErrNotFound           = errors.New("resource not found")
    ErrUserNotFound       = errors.New("user not found")
    ErrTaskNotFound       = errors.New("task not found")
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrEmailTaken         = errors.New("email already taken")
    ErrUsernameTaken      = errors.New("username already taken")
    ErrUnauthorized       = errors.New("unauthorized")
    ErrForbidden          = errors.New("forbidden")
    ErrInvalidToken       = errors.New("invalid token")
    ErrTokenExpired       = errors.New("token expired")
)
```

**Usage:**
```go
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
    user, err := s.userRepo.FindByEmail(ctx, req.Email)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, apperrors.ErrInvalidCredentials
    }
    // ...
}
```

---

## Error Checking Patterns

### Basic Check

```go
result, err := doSomething()
if err != nil {
    return err  // Or handle error
}
// Use result
```

### Checking for Specific Errors

```go
// Using errors.Is (Go 1.13+)
if errors.Is(err, apperrors.ErrNotFound) {
    // Handle not found
}

// Using errors.As for error types
var appErr *AppError
if errors.As(err, &appErr) {
    // Access appErr fields
    fmt.Println(appErr.Code)
}
```

### Early Return Pattern

```go
// Avoid deeply nested if-else
func processUser(id uint) error {
    user, err := findUser(id)
    if err != nil {
        return err
    }

    if !user.IsActive {
        return ErrUserInactive
    }

    if user.Role != "admin" {
        return ErrNotAuthorized
    }

    // Main logic here
    return nil
}
```

---

## Custom Error Types

### Rich Error Type

```go
type AppError struct {
    Code       string  // Machine-readable code
    Message    string  // Human-readable message
    StatusCode int     // HTTP status
    Err        error   // Wrapped error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

// Implement Unwrap for errors.Is/As
func (e *AppError) Unwrap() error {
    return e.Err
}
```

### Constructor Functions

```go
func NewNotFoundError(resource string) *AppError {
    return &AppError{
        Code:       "NOT_FOUND",
        Message:    fmt.Sprintf("%s not found", resource),
        StatusCode: http.StatusNotFound,
    }
}

func NewValidationError(message string) *AppError {
    return &AppError{
        Code:       "VALIDATION_ERROR",
        Message:    message,
        StatusCode: http.StatusBadRequest,
    }
}

func WrapError(err error, message string) *AppError {
    return &AppError{
        Code:       "INTERNAL_ERROR",
        Message:    message,
        StatusCode: http.StatusInternalServerError,
        Err:        err,
    }
}
```

---

## Error Wrapping

Go 1.13+ introduced error wrapping with `%w`:

```go
// Wrap error with context
err := doSomething()
if err != nil {
    return fmt.Errorf("failed to process user %d: %w", userID, err)
}

// Later, check for wrapped error
if errors.Is(err, sql.ErrNoRows) {
    // Handle not found
}
```

### Error Chain Example

```go
// Repository layer
func (r *repo) FindByID(ctx context.Context, id uint) (*User, error) {
    var user User
    err := r.db.First(&user, id).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil  // Not found is not an error here
    }
    if err != nil {
        return nil, fmt.Errorf("database query failed: %w", err)
    }
    return &user, nil
}

// Service layer
func (s *service) GetUser(ctx context.Context, id uint) (*User, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    if user == nil {
        return nil, apperrors.ErrUserNotFound
    }
    return user, nil
}

// Handler layer
func (h *handler) GetUser(c *gin.Context) {
    user, err := h.service.GetUser(ctx, id)
    if err != nil {
        utils.ErrorResponse(c, err)  // Converts to HTTP response
        return
    }
    utils.OKResponse(c, user)
}
```

---

## HTTP Error Responses

See `internal/utils/response.go`:

### Error Response Helper

```go
func ErrorResponse(c *gin.Context, err error) {
    // Check for known application errors
    switch {
    case errors.Is(err, apperrors.ErrNotFound),
         errors.Is(err, apperrors.ErrUserNotFound),
         errors.Is(err, apperrors.ErrTaskNotFound):
        NotFoundResponse(c, err.Error())

    case errors.Is(err, apperrors.ErrInvalidCredentials),
         errors.Is(err, apperrors.ErrUnauthorized),
         errors.Is(err, apperrors.ErrInvalidToken),
         errors.Is(err, apperrors.ErrTokenExpired):
        UnauthorizedResponse(c, err.Error())

    case errors.Is(err, apperrors.ErrForbidden):
        ForbiddenResponse(c, err.Error())

    case errors.Is(err, apperrors.ErrEmailTaken),
         errors.Is(err, apperrors.ErrUsernameTaken):
        ConflictResponse(c, err.Error())

    default:
        // Log unexpected errors
        logger.Error("Internal error", zap.Error(err))
        InternalErrorResponse(c, "An unexpected error occurred")
    }
}
```

### Standard Response Format

```go
type APIResponse[T any] struct {
    Success bool         `json:"success"`
    Data    T            `json:"data,omitempty"`
    Error   *ErrorDetail `json:"error,omitempty"`
    Meta    *Meta        `json:"meta,omitempty"`
}

type ErrorDetail struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details any    `json:"details,omitempty"`
}
```

### Response Helpers

```go
func OKResponse(c *gin.Context, data any) {
    c.JSON(http.StatusOK, APIResponse[any]{
        Success: true,
        Data:    data,
    })
}

func NotFoundResponse(c *gin.Context, message string) {
    c.JSON(http.StatusNotFound, APIResponse[any]{
        Success: false,
        Error: &ErrorDetail{
            Code:    "NOT_FOUND",
            Message: message,
        },
    })
}

func InternalErrorResponse(c *gin.Context, message string) {
    c.JSON(http.StatusInternalServerError, APIResponse[any]{
        Success: false,
        Error: &ErrorDetail{
            Code:    "INTERNAL_ERROR",
            Message: message,
        },
    })
}
```

---

## Panic and Recover

### Panic

Panic is for unrecoverable errors - use sparingly!

```go
// Panic stops normal execution
func MustGetConfig() *Config {
    cfg, err := LoadConfig()
    if err != nil {
        panic("failed to load config: " + err.Error())
    }
    return cfg
}

// Only use for:
// - Programming errors (should never happen)
// - Initialization failures
// - Assert-like checks in development
```

### Recover

Recover catches panics - usually in middleware:

```go
func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // Log the panic
                logger.Error("Panic recovered",
                    zap.Any("error", err),
                    zap.String("stack", string(debug.Stack())),
                )

                // Return 500 response
                c.JSON(http.StatusInternalServerError, APIResponse[any]{
                    Success: false,
                    Error: &ErrorDetail{
                        Code:    "INTERNAL_ERROR",
                        Message: "An unexpected error occurred",
                    },
                })
                c.Abort()
            }
        }()

        c.Next()
    }
}
```

---

## Best Practices

### DO:

```go
// 1. Handle errors at every level
result, err := doSomething()
if err != nil {
    return nil, err
}

// 2. Add context when wrapping
if err != nil {
    return fmt.Errorf("creating user %s: %w", email, err)
}

// 3. Use sentinel errors for known conditions
if user == nil {
    return apperrors.ErrUserNotFound
}

// 4. Log errors at the boundary (handler level)
func (h *handler) Create(c *gin.Context) {
    result, err := h.service.Create(ctx, req)
    if err != nil {
        logger.Error("Failed to create", zap.Error(err))
        utils.ErrorResponse(c, err)
        return
    }
}
```

### DON'T:

```go
// 1. Don't ignore errors
result, _ := doSomething()  // BAD!

// 2. Don't panic for recoverable errors
if user == nil {
    panic("user not found")  // BAD! Use error instead
}

// 3. Don't return generic errors everywhere
return errors.New("error")  // BAD! Not helpful

// 4. Don't log the same error multiple times
// Log once at the boundary, not in every layer
```

---

## Comparison: Node.js vs Go

### Node.js Exception Handling

```typescript
// Service
async findUser(id: number): Promise<User> {
  const user = await this.repo.findOne(id);
  if (!user) {
    throw new NotFoundException(`User ${id} not found`);
  }
  return user;
}

// Controller
@Get(':id')
async getUser(@Param('id') id: number) {
  try {
    return await this.userService.findUser(id);
  } catch (error) {
    if (error instanceof NotFoundException) {
      throw error;  // NestJS handles it
    }
    throw new InternalServerErrorException();
  }
}
```

### Go Error Handling

```go
// Service
func (s *service) FindUser(ctx context.Context, id uint) (*User, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("finding user: %w", err)
    }
    if user == nil {
        return nil, apperrors.ErrUserNotFound
    }
    return user, nil
}

// Handler
func (h *handler) GetUser(c *gin.Context) {
    user, err := h.service.FindUser(c.Request.Context(), id)
    if err != nil {
        utils.ErrorResponse(c, err)
        return
    }
    utils.OKResponse(c, user)
}
```

| Aspect | JavaScript | Go |
|--------|------------|-----|
| Mechanism | Exceptions (throw/catch) | Return values |
| Propagation | Automatic (bubbles up) | Manual (must return) |
| Handling | try/catch blocks | if err != nil |
| Visibility | Hidden control flow | Explicit in code |
| Performance | Stack unwinding | Just a value |

---

## Exercises

1. Look at `internal/errors/errors.go` - what sentinel errors are defined?
2. Look at `internal/utils/response.go` - how does ErrorResponse map errors to HTTP status?
3. Add a new error type `ErrRateLimited` and return it from rate limit middleware
4. Add error wrapping with context in the task service
