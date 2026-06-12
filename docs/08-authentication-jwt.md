# 08 - Authentication with JWT

JWT (JSON Web Token) authentication is stateless and perfect for APIs. Our implementation uses access + refresh token pattern.

## JWT Overview

### What is JWT?

A JWT consists of three parts:
```
header.payload.signature
```

**Example:**
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
eyJ1c2VyX2lkIjoxLCJleHAiOjE3MDk5MTIwMDB9.
9vZ5qGN8rS_5TvPKz9U5nYvFz2FwKvPvJq9W1rD3Q4o
```

- **Header**: Algorithm and token type
- **Payload**: Claims (user data, expiry, etc.)
- **Signature**: Verifies token wasn't tampered with

### Access + Refresh Token Pattern

```
┌─────────┐                    ┌─────────┐
│  Client │                    │  Server │
└────┬────┘                    └────┬────┘
     │                              │
     │ POST /login (credentials)    │
     │─────────────────────────────>│
     │                              │
     │ { accessToken, refreshToken }│
     │<─────────────────────────────│
     │                              │
     │ GET /tasks (accessToken)     │
     │─────────────────────────────>│
     │                              │
     │ { tasks }                    │
     │<─────────────────────────────│
     │                              │
     │ (accessToken expires)        │
     │                              │
     │ POST /refresh (refreshToken) │
     │─────────────────────────────>│
     │                              │
     │ { newAccessToken }           │
     │<─────────────────────────────│
```

| Token | Lifetime | Purpose |
|-------|----------|---------|
| Access Token | Short (15min-24h) | API authentication |
| Refresh Token | Long (7-30 days) | Get new access tokens |

---

## Our Implementation

### JWT Configuration

See `.env.example`:
```env
JWT_SECRET=your-super-secret-key-change-in-production
JWT_EXPIRY_HOURS=24
JWT_REFRESH_EXPIRY_HOURS=168
```

See `internal/config/config.go`:
```go
type JWTConfig struct {
    Secret             string        `mapstructure:"JWT_SECRET"`
    ExpiryHours        int           `mapstructure:"JWT_EXPIRY_HOURS"`
    RefreshExpiryHours int           `mapstructure:"JWT_REFRESH_EXPIRY_HOURS"`
    Expiry             time.Duration // Computed
    RefreshExpiry      time.Duration // Computed
}
```

### JWT Claims

See `internal/utils/jwt.go`:

**NestJS (passport-jwt):**
```typescript
interface JwtPayload {
  sub: number;     // user ID
  username: string;
  iat: number;     // issued at
  exp: number;     // expiry
}
```

**Go:**
```go
type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    TokenType string `json:"token_type"`  // "access" or "refresh"
    jwt.RegisteredClaims                   // Standard claims (exp, iat, etc.)
}
```

### Token Generation

```go
type JWTManager struct {
    secret        []byte
    expiry        time.Duration
    refreshExpiry time.Duration
}

func NewJWTManager(cfg *config.JWTConfig) *JWTManager {
    return &JWTManager{
        secret:        []byte(cfg.Secret),
        expiry:        cfg.Expiry,
        refreshExpiry: cfg.RefreshExpiry,
    }
}

func (m *JWTManager) GenerateAccessToken(userID uint, username string) (string, error) {
    claims := &Claims{
        UserID:    userID,
        Username:  username,
        TokenType: "access",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.expiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "task-manager",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(m.secret)
}
```

### Token Validation

```go
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(
        tokenString,
        &Claims{},
        func(token *jwt.Token) (interface{}, error) {
            // Verify signing method
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method")
            }
            return m.secret, nil
        },
    )

    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }

    return claims, nil
}
```

---

## Password Hashing

See `internal/utils/password.go`:

**Node.js (bcrypt):**
```javascript
const bcrypt = require('bcrypt');
const hash = await bcrypt.hash(password, 10);
const match = await bcrypt.compare(password, hash);
```

**Go (bcrypt):**
```go
import "golang.org/x/crypto/bcrypt"

// Hash password
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// Verify password
func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

**bcrypt.DefaultCost** is 10, which is a good balance between security and performance.

---

## Auth Service

See `internal/services/auth_service.go`:

### Register Flow

```go
func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
    // 1. Check if user exists
    existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
    if existing != nil {
        return nil, apperrors.ErrEmailTaken
    }

    // 2. Hash password
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }

    // 3. Create user
    user := &models.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: hashedPassword,
    }

    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    // 4. Generate tokens
    accessToken, _ := s.jwtManager.GenerateAccessToken(user.ID, user.Username)
    refreshToken, _ := s.jwtManager.GenerateRefreshToken(user.ID, user.Username)

    return &dto.AuthResponse{
        User:         user.ToResponse(),
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }, nil
}
```

### Login Flow

```go
func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
    // 1. Find user by email or username
    user, err := s.userRepo.FindByEmail(ctx, req.Login)
    if err != nil {
        return nil, err
    }
    if user == nil {
        user, _ = s.userRepo.FindByUsername(ctx, req.Login)
    }
    if user == nil {
        return nil, apperrors.ErrInvalidCredentials
    }

    // 2. Verify password
    if !utils.CheckPassword(req.Password, user.PasswordHash) {
        return nil, apperrors.ErrInvalidCredentials
    }

    // 3. Generate tokens
    accessToken, _ := s.jwtManager.GenerateAccessToken(user.ID, user.Username)
    refreshToken, _ := s.jwtManager.GenerateRefreshToken(user.ID, user.Username)

    return &dto.AuthResponse{
        User:         user.ToResponse(),
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }, nil
}
```

### Token Refresh Flow

```go
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenResponse, error) {
    // 1. Validate refresh token
    claims, err := s.jwtManager.ValidateToken(refreshToken)
    if err != nil {
        return nil, apperrors.ErrInvalidToken
    }

    // 2. Verify it's a refresh token
    if claims.TokenType != "refresh" {
        return nil, apperrors.ErrInvalidToken
    }

    // 3. Get user (ensure still exists)
    user, err := s.userRepo.FindByID(ctx, claims.UserID)
    if err != nil || user == nil {
        return nil, apperrors.ErrUserNotFound
    }

    // 4. Generate new access token
    newAccessToken, _ := s.jwtManager.GenerateAccessToken(user.ID, user.Username)

    return &dto.TokenResponse{
        AccessToken: newAccessToken,
        TokenType:   "Bearer",
        ExpiresIn:   int(s.jwtManager.GetExpiry().Seconds()),
    }, nil
}
```

---

## Auth Middleware

See `internal/middleware/auth.go`:

**Express.js (passport):**
```javascript
const passport = require('passport');
const JwtStrategy = require('passport-jwt').Strategy;

passport.use(new JwtStrategy(opts, (payload, done) => {
  User.findById(payload.sub, (err, user) => {
    if (user) return done(null, user);
    return done(null, false);
  });
}));

app.use('/protected', passport.authenticate('jwt', { session: false }));
```

**Gin:**
```go
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Get Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            utils.UnauthorizedResponse(c, "Authorization header required")
            c.Abort()
            return
        }

        // 2. Parse Bearer token
        if !strings.HasPrefix(authHeader, "Bearer ") {
            utils.UnauthorizedResponse(c, "Invalid authorization format")
            c.Abort()
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")

        // 3. Validate token
        claims, err := authService.ValidateAccessToken(token)
        if err != nil {
            utils.UnauthorizedResponse(c, "Invalid or expired token")
            c.Abort()
            return
        }

        // 4. Set user in context
        c.Set("userID", claims.UserID)
        c.Set("user", claims)

        c.Next()
    }
}
```

### Using Authenticated User

```go
// In protected handler
func (h *TaskHandler) Create(c *gin.Context) {
    // Get user ID from context (set by middleware)
    userID := middleware.MustGetUserID(c)

    var req dto.CreateTaskRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ValidationErrorResponse(c, err)
        return
    }

    // Create task for this user
    task, err := h.taskService.Create(c.Request.Context(), userID, &req)
    // ...
}
```

---

## Auth Handler

See `internal/handlers/auth_handler.go`:

```go
type AuthHandler struct {
    authService *services.AuthService
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req dto.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ValidationErrorResponse(c, err)
        return
    }

    response, err := h.authService.Register(c.Request.Context(), &req)
    if err != nil {
        utils.ErrorResponse(c, err)
        return
    }

    utils.CreatedResponse(c, response)
}

func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
    auth := rg.Group("/auth")
    {
        // Public routes
        auth.POST("/register", h.Register)
        auth.POST("/login", h.Login)
        auth.POST("/refresh", h.RefreshToken)

        // Protected routes
        auth.GET("/me", authMiddleware, h.Me)
    }
}
```

---

## Security Considerations

### DO:
- Use HTTPS in production
- Use strong secrets (32+ bytes, random)
- Keep access token lifetime short
- Validate token type (access vs refresh)
- Hash passwords with bcrypt

### DON'T:
- Store tokens in localStorage (use httpOnly cookies for web)
- Include sensitive data in JWT payload
- Use weak secrets or share them
- Skip token validation
- Use MD5/SHA1 for passwords

### Token Storage (Frontend)

| Storage | Pros | Cons |
|---------|------|------|
| localStorage | Easy to use | XSS vulnerable |
| sessionStorage | Tab-scoped | XSS vulnerable |
| httpOnly cookie | XSS safe | CSRF vulnerable |
| Memory | Most secure | Lost on refresh |

**Recommendation for API:** httpOnly cookies with CSRF protection, or memory + refresh token in httpOnly cookie.

---

## Testing Auth

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"john@example.com","password":"password123"}'

# Access protected route
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Refresh token
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"YOUR_REFRESH_TOKEN"}'
```

---

## Exercises

1. Look at `internal/utils/jwt.go` - what claims are included?
2. Look at `internal/middleware/auth.go` - what happens if token is missing?
3. Try to decode a JWT at jwt.io - what can you see in the payload?
4. Add a `Logout` endpoint that invalidates refresh tokens (hint: you'll need a token blacklist)
