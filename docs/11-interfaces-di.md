# 11 - Interfaces and Dependency Injection

Go's interfaces are implicit and powerful. They're central to writing testable, maintainable code.

## Interfaces in Go

### Implicit Implementation

In Go, you don't declare that a type implements an interface - it just does if it has the right methods.

**TypeScript (explicit):**
```typescript
interface UserRepository {
    findById(id: number): Promise<User | null>;
    create(user: User): Promise<User>;
}

// Must explicitly implement
class MySQLUserRepository implements UserRepository {
    async findById(id: number): Promise<User | null> { ... }
    async create(user: User): Promise<User> { ... }
}
```

**Go (implicit):**
```go
type UserRepository interface {
    FindByID(ctx context.Context, id uint) (*User, error)
    Create(ctx context.Context, user *User) error
}

// No "implements" keyword!
// If it has these methods, it implements the interface
type mysqlUserRepository struct {
    db *gorm.DB
}

func (r *mysqlUserRepository) FindByID(ctx context.Context, id uint) (*User, error) {
    // ...
}

func (r *mysqlUserRepository) Create(ctx context.Context, user *User) error {
    // ...
}
```

### Interface Satisfaction

```go
// Check at compile time
var _ UserRepository = (*mysqlUserRepository)(nil)

// This line will fail to compile if mysqlUserRepository
// doesn't implement all UserRepository methods
```

---

## Our Repository Pattern

See `internal/repository/repository.go`:

```go
// Interface definition
type UserRepository interface {
    FindByID(ctx context.Context, id uint) (*models.User, error)
    FindByEmail(ctx context.Context, email string) (*models.User, error)
    Create(ctx context.Context, user *models.User) error
    Update(ctx context.Context, user *models.User) error
    Delete(ctx context.Context, id uint) error
}

// Implementation
type userRepository struct {
    db *gorm.DB
}

// Constructor returns the interface, not the concrete type
func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}
```

### Why This Pattern?

1. **Testability**: Can create mock implementations
2. **Flexibility**: Can swap implementations (MySQL → PostgreSQL)
3. **Abstraction**: Service doesn't know about GORM

---

## Dependency Injection

### What is DI?

Dependency Injection means passing dependencies to a component instead of having it create them.

**Without DI (bad):**
```go
type AuthService struct {}

func (s *AuthService) Login(email, password string) {
    // Creates its own repository - hard to test!
    repo := NewUserRepository(getGlobalDB())
    user := repo.FindByEmail(email)
}
```

**With DI (good):**
```go
type AuthService struct {
    userRepo UserRepository  // Interface, not concrete type
}

func NewAuthService(userRepo UserRepository) *AuthService {
    return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Login(email, password string) {
    // Uses injected repository
    user := s.userRepo.FindByEmail(email)
}
```

### DI in Our Application

See `cmd/api/main.go`:

```go
// Create repositories
userRepo := repository.NewUserRepository(db)
taskRepo := repository.NewTaskRepository(db)

// Inject into services
authSvc := services.NewAuthService(userRepo, &cfg.JWT, notificationSvc)
taskSvc := services.NewTaskService(taskRepo, notificationSvc)

// Inject into handlers
authHandler := handlers.NewAuthHandler(authSvc)
```

---

## Comparison with NestJS

**NestJS DI:**
```typescript
@Injectable()
export class AuthService {
    constructor(
        @InjectRepository(User)
        private userRepository: Repository<User>,
        private jwtService: JwtService,
    ) {}
}
```

**Go DI:**
```go
type AuthService struct {
    userRepo   UserRepository
    jwtManager *JWTManager
}

func NewAuthService(userRepo UserRepository, jwt *JWTManager) *AuthService {
    return &AuthService{
        userRepo:   userRepo,
        jwtManager: jwt,
    }
}
```

| Aspect | NestJS | Go |
|--------|--------|-----|
| DI Container | Built-in (IoC container) | Manual (or use wire/fx) |
| Registration | `@Injectable()` decorator | Constructor functions |
| Injection | `@Inject()` decorators | Explicit parameter passing |
| Lifetime | Singleton by default | You decide |

---

## Testing with Interfaces

### Creating Mocks

```go
// Mock implementation for testing
type mockUserRepository struct {
    users map[uint]*models.User
}

func (m *mockUserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
    user, ok := m.users[id]
    if !ok {
        return nil, nil
    }
    return user, nil
}

func (m *mockUserRepository) Create(ctx context.Context, user *models.User) error {
    user.ID = uint(len(m.users) + 1)
    m.users[user.ID] = user
    return nil
}
// ... implement other methods
```

### Using in Tests

```go
func TestAuthService_Login(t *testing.T) {
    // Create mock
    mockRepo := &mockUserRepository{
        users: map[uint]*models.User{
            1: {
                ID:           1,
                Email:        "test@example.com",
                PasswordHash: hashPassword("password123"),
            },
        },
    }

    // Inject mock
    authService := NewAuthService(mockRepo, jwtConfig)

    // Test
    result, err := authService.Login(context.Background(), &LoginRequest{
        Login:    "test@example.com",
        Password: "password123",
    })

    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

---

## Interface Design Guidelines

### Small Interfaces

Go convention: Keep interfaces small.

```go
// Good - single responsibility
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// Compose when needed
type ReadWriter interface {
    Reader
    Writer
}
```

### Accept Interfaces, Return Structs

```go
// Good - accept interface
func ProcessUser(repo UserRepository) { ... }

// Good - return concrete type
func NewUserRepository(db *gorm.DB) *userRepository { ... }
// Or return interface if you need to hide implementation
func NewUserRepository(db *gorm.DB) UserRepository { ... }
```

### Define Interfaces Near Usage

```go
// In the consumer package, not the provider
package service

type UserRepository interface {
    FindByID(ctx context.Context, id uint) (*User, error)
}

type AuthService struct {
    userRepo UserRepository
}
```

---

## Common Patterns

### The Empty Interface

```go
// interface{} accepts any type (like TypeScript's `any`)
func PrintAnything(v interface{}) {
    fmt.Println(v)
}

// Go 1.18+ has `any` as an alias
func PrintAnything(v any) {
    fmt.Println(v)
}
```

### Type Assertions

```go
var i interface{} = "hello"

// Type assertion
s := i.(string)  // Panics if wrong type

// Safe type assertion
s, ok := i.(string)
if ok {
    fmt.Println(s)
}

// Type switch
switch v := i.(type) {
case string:
    fmt.Println("String:", v)
case int:
    fmt.Println("Int:", v)
default:
    fmt.Println("Unknown type")
}
```

---

## Exercises

1. Look at `internal/repository/repository.go` - how many interfaces are defined?
2. Trace the DI chain in `main.go` from database to handler
3. Create a mock `TaskRepository` for testing
4. Try adding a new method to an interface and see what breaks
