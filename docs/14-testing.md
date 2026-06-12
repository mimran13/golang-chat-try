# 14 - Testing in Go

Go has excellent built-in testing support. No external frameworks required (though we use testify for assertions).

## Go Testing vs JavaScript Testing

| Aspect | Go | JavaScript |
|--------|-----|------------|
| Test Runner | Built-in (`go test`) | Jest, Mocha, Vitest |
| Assertions | testify/assert | Jest expect, Chai |
| Mocking | testify/mock, gomock | Jest mocks, Sinon |
| Coverage | Built-in (`-cover`) | c8, istanbul |
| Benchmarks | Built-in | Benchmark.js |

---

## Test File Conventions

```
internal/services/
├── auth_service.go       # Implementation
├── auth_service_test.go  # Tests (same package)
```

- Test files end with `_test.go`
- Test functions start with `Test`
- Can be in same package (white-box) or `_test` package (black-box)

---

## Writing Tests

### Basic Test

```go
// math_test.go
package math

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}
```

### Using testify/assert

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
    result := Add(2, 3)

    // assert continues on failure
    assert.Equal(t, 5, result)
    assert.NotNil(t, result)
    assert.True(t, result > 0)

    // require stops on failure
    require.NoError(t, err)
    require.NotNil(t, user)
}
```

### Table-Driven Tests

The Go idiom for testing multiple cases:

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -1, -2, -3},
        {"zero", 0, 0, 0},
        {"mixed", -1, 5, 4},
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            result := Add(tc.a, tc.b)
            assert.Equal(t, tc.expected, result)
        })
    }
}
```

**JavaScript equivalent:**
```javascript
describe('Add', () => {
    test.each([
        [2, 3, 5],
        [-1, -2, -3],
        [0, 0, 0],
    ])('add(%i, %i) = %i', (a, b, expected) => {
        expect(add(a, b)).toBe(expected);
    });
});
```

---

## Our Test Structure

```
internal/
├── testutil/
│   ├── helpers.go      # Test utilities
│   ├── mocks.go        # Mock implementations
│   └── integration.go  # Integration test setup
├── utils/
│   ├── jwt_test.go         # Unit tests
│   └── password_test.go    # Unit tests
├── services/
│   ├── auth_service_test.go    # Unit tests with mocks
│   └── task_service_test.go    # Unit tests with mocks
├── middleware/
│   ├── auth_test.go        # Middleware tests
│   └── ratelimit_test.go   # Middleware tests
└── handlers/
    ├── auth_handler_test.go   # Integration tests
    └── task_handler_test.go   # Integration tests
```

---

## Unit Tests with Mocks

### Creating Mocks

```go
// testutil/mocks.go
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.User), args.Error(1)
}
```

### Using Mocks in Tests

```go
func TestAuthService_Login(t *testing.T) {
    // Create mock
    mockRepo := new(testutil.MockUserRepository)

    // Set expectations
    mockRepo.On("FindByEmailOrUsername", ctx, "test@example.com").
        Return(user, nil)

    // Create service with mock
    service := NewAuthService(mockRepo, jwtConfig)

    // Test
    result, err := service.Login(ctx, &LoginRequest{
        Login:    "test@example.com",
        Password: "password123",
    })

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)

    // Verify mock was called
    mockRepo.AssertExpectations(t)
}
```

---

## Integration Tests

### Setup Test Application

```go
// testutil/integration.go
func SetupTestApp(t *testing.T) *TestApp {
    // Use SQLite for fast tests
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&models.User{}, &models.Task{})

    // Create real dependencies
    userRepo := repository.NewUserRepository(db)
    authSvc := services.NewAuthService(userRepo, jwtConfig)

    // Create router
    router := gin.New()
    // ... setup routes

    return &TestApp{DB: db, Router: router}
}
```

### Writing Integration Tests

```go
func TestAuthHandler_Register_Integration(t *testing.T) {
    app := testutil.SetupTestApp(t)
    defer app.Cleanup()

    body := map[string]interface{}{
        "username": "newuser",
        "email":    "new@example.com",
        "password": "password123",
    }

    req := testutil.MakeRequest("POST", "/api/v1/auth/register", body)
    rec := testutil.ExecuteRequest(app.Router, req)

    assert.Equal(t, http.StatusCreated, rec.Code)

    var response dto.APIResponse[dto.AuthResponse]
    json.Unmarshal(rec.Body.Bytes(), &response)
    assert.True(t, response.Success)
    assert.Equal(t, "newuser", response.Data.User.Username)
}
```

---

## Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package
go test ./internal/services/...

# Run specific test
go test -run TestAuthService_Login ./internal/services/

# Skip long tests
go test -short ./...

# With coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run with race detector
go test -race ./...

# Run benchmarks
go test -bench=. ./...
```

---

## Test Helpers

### Test Fixtures

```go
// testutil/helpers.go
func TestUser() *models.User {
    return &models.User{
        Username:     "testuser",
        Email:        "test@example.com",
        PasswordHash: "hashed_password",
    }
}

func TestTask(userID uint) *models.Task {
    return &models.Task{
        UserID: userID,
        Title:  "Test Task",
        Status: models.TaskStatusPending,
    }
}
```

### HTTP Request Helpers

```go
func MakeRequest(method, url string, body interface{}) *http.Request {
    jsonBytes, _ := json.Marshal(body)
    req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
    req.Header.Set("Content-Type", "application/json")
    return req
}

func MakeAuthRequest(method, url string, body interface{}, token string) *http.Request {
    req := MakeRequest(method, url, body)
    req.Header.Set("Authorization", "Bearer "+token)
    return req
}

func ExecuteRequest(router *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
    rec := httptest.NewRecorder()
    router.ServeHTTP(rec, req)
    return rec
}
```

---

## Testing Patterns

### Testing Errors

```go
func TestLogin_InvalidCredentials(t *testing.T) {
    mockRepo := new(MockUserRepository)
    mockRepo.On("FindByEmail", ctx, "test@example.com").Return(nil, nil)

    service := NewAuthService(mockRepo, jwtConfig)
    _, err := service.Login(ctx, &LoginRequest{...})

    // Check for specific error
    assert.ErrorIs(t, err, apperrors.ErrInvalidCredentials)
}
```

### Testing Async Operations

```go
func TestNotificationService(t *testing.T) {
    svc := NewNotificationService(3, 10)
    defer svc.Shutdown()

    // Send notification
    svc.SendWelcomeEmail("test@example.com", "testuser")

    // Wait for processing
    time.Sleep(100 * time.Millisecond)

    // Check stats
    stats := svc.GetStats()
    assert.Equal(t, 1, stats.Processed)
}
```

### Testing with Context Timeout

```go
func TestWithTimeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()

    // Operation that might timeout
    _, err := slowOperation(ctx)

    if errors.Is(err, context.DeadlineExceeded) {
        t.Log("Operation timed out as expected")
    }
}
```

---

## Benchmarks

```go
func BenchmarkHashPassword(b *testing.B) {
    password := "password123"

    for i := 0; i < b.N; i++ {
        HashPassword(password)
    }
}

func BenchmarkValidateToken(b *testing.B) {
    manager := NewJWTManager(config)
    token, _ := manager.GenerateAccessToken(1, "user")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        manager.ValidateToken(token)
    }
}
```

**Running benchmarks:**
```bash
go test -bench=. -benchmem ./internal/utils/

# Output:
# BenchmarkHashPassword-8        10    108433275 ns/op    5024 B/op    10 allocs/op
# BenchmarkValidateToken-8   500000         2541 ns/op     832 B/op    12 allocs/op
```

---

## Makefile Commands

```makefile
## test: Run all tests
test:
	go test -v ./...

## test-coverage: Run tests with coverage report
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

## test-race: Run tests with race detector
test-race:
	go test -race ./...

## test-short: Run only short tests
test-short:
	go test -short ./...
```

---

## Exercises

1. Run `make test` and check that all tests pass
2. Run `make test-coverage` and review the coverage report
3. Add a test for a new validation rule
4. Write a benchmark for the rate limiter
5. Add an integration test for updating task status
