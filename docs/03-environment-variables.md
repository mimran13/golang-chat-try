# 03 - Environment Variables

## Reading Environment Variables

### Basic Access

**Node.js:**
```javascript
const port = process.env.PORT;
const nodeEnv = process.env.NODE_ENV || 'development';
```

**Go:**
```go
import "os"

port := os.Getenv("PORT")              // Returns "" if not set
nodeEnv := os.Getenv("NODE_ENV")       // Returns "" if not set

// Check if variable exists
value, exists := os.LookupEnv("PORT")  // exists is bool
if !exists {
    value = "8080"  // default
}
```

### Key Difference

| Aspect | Node.js | Go |
|--------|---------|-----|
| Access | `process.env.KEY` | `os.Getenv("KEY")` |
| Missing value | `undefined` | `""` (empty string) |
| Check existence | `KEY in process.env` | `os.LookupEnv("KEY")` returns `(value, bool)` |
| Global object | Yes (`process.env`) | No (use functions) |

---

## Loading .env Files

### Using godotenv

In our project, we use `godotenv`, similar to the `dotenv` package in Node.js:

**Node.js:**
```javascript
import dotenv from 'dotenv';
dotenv.config();
// Now process.env has values from .env
```

**Go:**
```go
import "github.com/joho/godotenv"

// Load .env file into environment
err := godotenv.Load()
if err != nil {
    log.Println("No .env file found")
}
// Now os.Getenv() returns values from .env
```

### Load Options

```go
// Load specific file
godotenv.Load(".env.local")

// Load multiple files (later files override)
godotenv.Load(".env", ".env.local")

// Load into a map instead of environment
envMap, err := godotenv.Read(".env")
fmt.Println(envMap["PORT"])  // Access directly from map
```

---

## Configuration Patterns

### Pattern 1: Direct Access (Simple but Scattered)

```go
// Not recommended for larger apps - env access scattered everywhere
func main() {
    port := os.Getenv("PORT")
    dbHost := os.Getenv("DB_HOST")
    // ...
}
```

### Pattern 2: Centralized Config Struct (Our Approach)

```go
// config/config.go
type Config struct {
    Port        int
    Environment string
    DBHost      string
}

func Load() (*Config, error) {
    godotenv.Load()
    return &Config{
        Port:        getEnvAsInt("PORT", 8080),
        Environment: getEnv("ENVIRONMENT", "development"),
        DBHost:      getEnv("DB_HOST", "localhost"),
    }, nil
}

// main.go
func main() {
    cfg, _ := config.Load()
    fmt.Printf("Port: %d\n", cfg.Port)
}
```

**Benefits:**
- All config in one place
- Type safety (Port is `int`, not `string`)
- Easy to test (inject mock config)
- Clear documentation of required variables

### Pattern 3: Using viper (Advanced)

For complex apps, consider [viper](https://github.com/spf13/viper):

```go
import "github.com/spf13/viper"

viper.SetConfigFile(".env")
viper.ReadInConfig()

port := viper.GetInt("PORT")
dbHost := viper.GetString("DB_HOST")
```

Viper supports: env files, JSON, YAML, watching for changes, and more.

---

## Type Conversion

Environment variables are always strings. You need to convert them:

### String to Integer

```go
import "strconv"

portStr := os.Getenv("PORT")  // "8080"
port, err := strconv.Atoi(portStr)  // Atoi = "ASCII to Integer"
if err != nil {
    port = 8080  // default
}
```

### String to Boolean

```go
debugStr := os.Getenv("DEBUG")  // "true"
debug, err := strconv.ParseBool(debugStr)
// Accepts: "1", "t", "T", "TRUE", "true", "True"
// And: "0", "f", "F", "FALSE", "false", "False"
```

### Helper Functions

In our `config/config.go`, we have helpers:

```go
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    strValue := os.Getenv(key)
    if value, err := strconv.Atoi(strValue); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
    strValue := os.Getenv(key)
    if value, err := strconv.ParseBool(strValue); err == nil {
        return value
    }
    return defaultValue
}
```

---

## Comparison with NestJS ConfigModule

**NestJS:**
```typescript
// app.module.ts
@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      validationSchema: Joi.object({
        PORT: Joi.number().default(3000),
        NODE_ENV: Joi.string().default('development'),
      }),
    }),
  ],
})
export class AppModule {}

// Using in a service
@Injectable()
export class AppService {
  constructor(private configService: ConfigService) {}

  getPort(): number {
    return this.configService.get<number>('PORT');
  }
}
```

**Go (our approach):**
```go
// pkg/config/config.go
type Config struct {
    Port        int
    Environment string
}

func Load() (*Config, error) {
    // Load and validate
    return &Config{
        Port:        getEnvAsInt("PORT", 8080),
        Environment: getEnv("ENVIRONMENT", "development"),
    }, nil
}

// cmd/api/main.go
func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }
    // Pass cfg to handlers/services that need it
}
```

**Key Difference:** Go doesn't have dependency injection built-in. We pass config explicitly to functions/structs that need it.

---

## Best Practices

### 1. Never Hardcode Secrets

```go
// BAD
dbPassword := "secret123"

// GOOD
dbPassword := os.Getenv("DB_PASSWORD")
```

### 2. Always Have Defaults for Non-Sensitive Values

```go
// Port can have default
port := getEnvAsInt("PORT", 8080)

// But secrets should error if missing
dbPassword := os.Getenv("DB_PASSWORD")
if dbPassword == "" {
    log.Fatal("DB_PASSWORD is required")
}
```

### 3. Document Required Variables

Keep `.env.example` updated with all variables (without real secrets).

### 4. Validate at Startup

```go
func Load() (*Config, error) {
    cfg := &Config{}

    // Required variable
    cfg.DBPassword = os.Getenv("DB_PASSWORD")
    if cfg.DBPassword == "" {
        return nil, fmt.Errorf("DB_PASSWORD is required")
    }

    return cfg, nil
}
```

---

## Exercises

1. Add a new environment variable `LOG_LEVEL` with default "info"
2. Create a `getEnvAsBool` helper function
3. Try running the app without a `.env` file (defaults should work)
4. Add validation that fails if a required variable is missing

---

## Next Steps

- Read **04-http-server-basics.md** to learn about building APIs
- Try adding database configuration variables (we'll use them later)
