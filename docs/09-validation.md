# 09 - Request Validation

Go-playground/validator is integrated with Gin for request validation. It's similar to class-validator in NestJS.

## Overview

| Go | NestJS |
|----|--------|
| go-playground/validator | class-validator |
| `binding:"required"` | `@IsNotEmpty()` |
| `binding:"email"` | `@IsEmail()` |
| `binding:"min=8"` | `@MinLength(8)` |
| Struct tags | Decorators |

---

## Basic Validation

### Validation Tags

**NestJS (class-validator):**
```typescript
import { IsEmail, IsNotEmpty, MinLength } from 'class-validator';

export class RegisterDto {
  @IsNotEmpty()
  username: string;

  @IsEmail()
  email: string;

  @MinLength(8)
  password: string;
}
```

**Go (validator):**
```go
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}
```

### Common Validation Tags

| Tag | Description | Example |
|-----|-------------|---------|
| `required` | Field must be present and non-zero | `binding:"required"` |
| `email` | Valid email format | `binding:"email"` |
| `min=N` | Minimum length/value | `binding:"min=8"` |
| `max=N` | Maximum length/value | `binding:"max=100"` |
| `len=N` | Exact length | `binding:"len=10"` |
| `eq=X` | Equals value | `binding:"eq=active"` |
| `ne=X` | Not equals value | `binding:"ne=deleted"` |
| `oneof=A B C` | One of the values | `binding:"oneof=low medium high"` |
| `gt=N` | Greater than | `binding:"gt=0"` |
| `gte=N` | Greater than or equal | `binding:"gte=1"` |
| `lt=N` | Less than | `binding:"lt=100"` |
| `lte=N` | Less than or equal | `binding:"lte=99"` |
| `url` | Valid URL | `binding:"url"` |
| `uuid` | Valid UUID | `binding:"uuid"` |
| `numeric` | Numeric string | `binding:"numeric"` |
| `alpha` | Alphabetic only | `binding:"alpha"` |
| `alphanum` | Alphanumeric | `binding:"alphanum"` |

---

## Our DTOs

### Auth DTOs

See `internal/dto/auth.go`:

```go
// Registration request
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50,alphanum"`
    Email    string `json:"email" binding:"required,email,max=100"`
    Password string `json:"password" binding:"required,min=8,max=72"`
}

// Login request
type LoginRequest struct {
    Login    string `json:"login" binding:"required"`   // Email or username
    Password string `json:"password" binding:"required"`
}

// Refresh token request
type RefreshTokenRequest struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}
```

### Task DTOs

See `internal/dto/task.go`:

```go
type CreateTaskRequest struct {
    Title       string  `json:"title" binding:"required,min=1,max=200"`
    Description string  `json:"description" binding:"max=2000"`
    Priority    string  `json:"priority" binding:"omitempty,oneof=low medium high"`
    DueDate     *string `json:"due_date" binding:"omitempty"`
}

type UpdateTaskRequest struct {
    Title       *string `json:"title" binding:"omitempty,min=1,max=200"`
    Description *string `json:"description" binding:"omitempty,max=2000"`
    Status      *string `json:"status" binding:"omitempty,oneof=pending in_progress completed cancelled"`
    Priority    *string `json:"priority" binding:"omitempty,oneof=low medium high"`
    DueDate     *string `json:"due_date" binding:"omitempty"`
}
```

**Note:** Using `*string` (pointer) allows distinguishing between "not provided" and "empty string".

---

## Validation in Handlers

### Using ShouldBindJSON

```go
func (h *AuthHandler) Register(c *gin.Context) {
    var req dto.RegisterRequest

    // ShouldBindJSON:
    // 1. Reads JSON body
    // 2. Unmarshals into struct
    // 3. Validates using binding tags
    // 4. Returns error if validation fails
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ValidationErrorResponse(c, err)
        return
    }

    // req is now validated and ready to use
    response, err := h.authService.Register(c.Request.Context(), &req)
    // ...
}
```

### Binding Methods

```go
// Bind JSON body
c.ShouldBindJSON(&req)      // Returns error
c.BindJSON(&req)            // Aborts on error (avoid)

// Bind query parameters
c.ShouldBindQuery(&req)

// Bind form data
c.ShouldBind(&req)

// Bind URI parameters
c.ShouldBindUri(&req)
```

---

## Custom Validation Errors

### Default Error Format

Validator returns errors like:
```
Key: 'RegisterRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag
```

Not very user-friendly!

### Custom Error Response

See `internal/utils/response.go`:

```go
func ValidationErrorResponse(c *gin.Context, err error) {
    var ve validator.ValidationErrors
    if errors.As(err, &ve) {
        errors := make([]ValidationError, len(ve))
        for i, fe := range ve {
            errors[i] = ValidationError{
                Field:   toSnakeCase(fe.Field()),
                Message: getErrorMessage(fe),
            }
        }
        c.JSON(http.StatusBadRequest, APIResponse[any]{
            Success: false,
            Error: &ErrorDetail{
                Code:    "VALIDATION_ERROR",
                Message: "Validation failed",
                Details: errors,
            },
        })
        return
    }

    // JSON parsing error
    c.JSON(http.StatusBadRequest, APIResponse[any]{
        Success: false,
        Error: &ErrorDetail{
            Code:    "INVALID_JSON",
            Message: "Invalid JSON format",
        },
    })
}

func getErrorMessage(fe validator.FieldError) string {
    switch fe.Tag() {
    case "required":
        return "This field is required"
    case "email":
        return "Must be a valid email address"
    case "min":
        return fmt.Sprintf("Must be at least %s characters", fe.Param())
    case "max":
        return fmt.Sprintf("Must be at most %s characters", fe.Param())
    case "oneof":
        return fmt.Sprintf("Must be one of: %s", fe.Param())
    default:
        return fmt.Sprintf("Failed validation: %s", fe.Tag())
    }
}
```

### Result

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "email",
        "message": "Must be a valid email address"
      },
      {
        "field": "password",
        "message": "Must be at least 8 characters"
      }
    ]
  }
}
```

---

## Advanced Validation

### Custom Validators

```go
import "github.com/go-playground/validator/v10"

// Register custom validator
func setupValidators(v *validator.Validate) {
    v.RegisterValidation("password", validatePassword)
}

// Custom password validator
func validatePassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()

    // Check minimum length
    if len(password) < 8 {
        return false
    }

    // Check for at least one uppercase
    hasUpper := false
    hasLower := false
    hasNumber := false

    for _, c := range password {
        if unicode.IsUpper(c) {
            hasUpper = true
        }
        if unicode.IsLower(c) {
            hasLower = true
        }
        if unicode.IsNumber(c) {
            hasNumber = true
        }
    }

    return hasUpper && hasLower && hasNumber
}

// Usage in struct
type RegisterRequest struct {
    Password string `json:"password" binding:"required,password"`
}
```

### Conditional Validation

```go
// Skip validation if field is empty
type UpdateTaskRequest struct {
    Title *string `json:"title" binding:"omitempty,min=1,max=200"`
    // "omitempty" = don't validate if nil or zero value
}

// Required only if another field is present
type Request struct {
    Type    string `json:"type" binding:"required,oneof=email phone"`
    Email   string `json:"email" binding:"required_if=Type email,omitempty,email"`
    Phone   string `json:"phone" binding:"required_if=Type phone,omitempty,e164"`
}

// Required unless another field is present
type Request struct {
    Primary   string `json:"primary" binding:"required_without=Secondary"`
    Secondary string `json:"secondary" binding:"required_without=Primary"`
}
```

### Cross-Field Validation

```go
type DateRange struct {
    StartDate time.Time `json:"start_date" binding:"required"`
    EndDate   time.Time `json:"end_date" binding:"required,gtfield=StartDate"`
    // "gtfield=StartDate" = must be greater than StartDate
}

type Password struct {
    Password        string `json:"password" binding:"required,min=8"`
    ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
    // "eqfield=Password" = must equal Password field
}
```

---

## Query Parameter Validation

```go
type TaskListQuery struct {
    Page     int    `form:"page" binding:"omitempty,min=1"`
    Limit    int    `form:"limit" binding:"omitempty,min=1,max=100"`
    Status   string `form:"status" binding:"omitempty,oneof=pending in_progress completed cancelled"`
    Priority string `form:"priority" binding:"omitempty,oneof=low medium high"`
    SortBy   string `form:"sort_by" binding:"omitempty,oneof=created_at updated_at due_date title"`
    Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
}

func (h *TaskHandler) List(c *gin.Context) {
    var query TaskListQuery

    // Use ShouldBindQuery for query parameters
    if err := c.ShouldBindQuery(&query); err != nil {
        utils.ValidationErrorResponse(c, err)
        return
    }

    // Set defaults
    if query.Page == 0 {
        query.Page = 1
    }
    if query.Limit == 0 {
        query.Limit = 10
    }

    // Use validated query
    tasks, total, err := h.taskService.List(c.Request.Context(), userID, &query)
    // ...
}
```

---

## URI Parameter Validation

```go
type TaskURI struct {
    ID uint `uri:"id" binding:"required,min=1"`
}

func (h *TaskHandler) GetByID(c *gin.Context) {
    var uri TaskURI

    if err := c.ShouldBindUri(&uri); err != nil {
        utils.ValidationErrorResponse(c, err)
        return
    }

    task, err := h.taskService.GetByID(c.Request.Context(), userID, uri.ID)
    // ...
}
```

---

## Comparison: NestJS vs Go

### NestJS Validation Pipe

```typescript
// main.ts
app.useGlobalPipes(new ValidationPipe({
  whitelist: true,
  forbidNonWhitelisted: true,
  transform: true,
}));

// dto
export class CreateTaskDto {
  @IsString()
  @IsNotEmpty()
  @MaxLength(200)
  title: string;

  @IsString()
  @IsOptional()
  @MaxLength(2000)
  description?: string;

  @IsEnum(Priority)
  @IsOptional()
  priority?: Priority;
}
```

### Go Validation

```go
// dto
type CreateTaskRequest struct {
    Title       string  `json:"title" binding:"required,min=1,max=200"`
    Description string  `json:"description" binding:"omitempty,max=2000"`
    Priority    string  `json:"priority" binding:"omitempty,oneof=low medium high"`
}

// handler
func (h *Handler) CreateTask(c *gin.Context) {
    var req CreateTaskRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ValidationErrorResponse(c, err)
        return
    }
    // Use req
}
```

---

## Exercises

1. Look at `internal/dto/auth.go` - what validations are applied to the register request?
2. Look at `internal/dto/task.go` - why does UpdateTaskRequest use pointer types?
3. Add a custom validation to ensure due_date is in the future
4. Add password strength validation (uppercase, lowercase, number required)
