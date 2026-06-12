# 02 - Packages and Modules

## Understanding Packages

### What is a Package?

In Go, a **package** is a collection of `.go` files in the same directory. All files in a directory must have the same package name.

**Node.js (modules):**
```typescript
// userService.ts
export function getUser() { ... }
export const defaultUser = { ... }

// Using it:
import { getUser, defaultUser } from './userService'
```

**Go (packages):**
```go
// handlers/hello.go
package handlers

func HelloHandler() { ... }  // Exported (public)
func helper() { ... }        // Unexported (private)

// Using it in another package:
import "github.com/user/task-manager/internal/handlers"
handlers.HelloHandler()  // Works
handlers.helper()        // ERROR: cannot access unexported function
```

### Key Differences

| Concept | Node.js/TypeScript | Go |
|---------|-------------------|-----|
| Unit of code sharing | File (module) | Directory (package) |
| Export mechanism | `export` keyword | Capital first letter |
| Import mechanism | `import { x } from './file'` | `import "path/to/package"` |
| Private | No `export` / `private` keyword | Lowercase first letter |
| Public | `export` / `public` keyword | Uppercase first letter |

---

## Visibility Rules

### The Capital Letter Rule

This is one of Go's most distinctive features. **No keywords needed** for visibility:

```go
package handlers

// PUBLIC (Exported) - Capital first letter
func HelloHandler() { ... }      // Can be called from other packages
type HelloResponse struct { ... } // Can be used from other packages
var GlobalConfig = "..."          // Can be accessed from other packages

// PRIVATE (Unexported) - Lowercase first letter
func helper() { ... }            // Only usable within 'handlers' package
type internalData struct { ... } // Only usable within 'handlers' package
var secretKey = "..."            // Only usable within 'handlers' package
```

**Why?** It makes code self-documenting. Glancing at any identifier tells you its visibility.

### Struct Field Visibility

The same rules apply to struct fields:

```go
type User struct {
    ID        int    // Public - other packages can read/write
    Name      string // Public
    password  string // Private - only this package can access
}
```

**TypeScript equivalent:**
```typescript
class User {
    public id: number;
    public name: string;
    private password: string;
}
```

---

## Import Statements

### Basic Imports

```go
import "fmt"                    // Standard library
import "net/http"               // Standard library (nested)
import "github.com/joho/godotenv" // Third-party
import "github.com/user/task-manager/pkg/config" // Local package
```

### Grouped Imports (Preferred Style)

```go
import (
    // Standard library (first group)
    "fmt"
    "net/http"
    "time"

    // Third-party packages (second group)
    "github.com/joho/godotenv"

    // Local packages (third group)
    "github.com/user/task-manager/internal/handlers"
    "github.com/user/task-manager/pkg/config"
)
```

**Note:** `goimports` (included in golangci-lint) automatically organizes imports for you!

### Import Aliases

```go
import (
    // Alias for long package names
    cfg "github.com/user/task-manager/pkg/config"

    // Alias to avoid conflicts
    gohttp "net/http"
    customhttp "myproject/http"
)

// Using aliases:
cfg.Load()
gohttp.ListenAndServe(...)
```

### Special Imports

```go
import (
    // Blank import - runs package's init() but doesn't use package
    // Common for database drivers that register themselves
    _ "github.com/go-sql-driver/mysql"

    // Dot import - imports all exported names into current namespace
    // Generally discouraged (pollutes namespace)
    . "fmt"
)

// With dot import:
Println("Hello") // Instead of fmt.Println()
```

---

## Package Initialization

### The init() Function

Go packages can have an `init()` function that runs automatically when the package is imported:

```go
package config

var settings map[string]string

// init runs BEFORE main(), when package is first imported
func init() {
    settings = make(map[string]string)
    settings["default_port"] = "8080"
    fmt.Println("Config package initialized")
}
```

**Node.js equivalent:**
```typescript
// This runs when file is imported
const settings = new Map();
settings.set("default_port", "8080");
console.log("Config module loaded");
```

**Key differences:**
- Multiple `init()` functions allowed per package (all run)
- Run order: imported packages' `init()` → current package's `init()` → `main()`
- Runs once per program, even if imported multiple times

---

## The main Package

### Special Rules for main

The `main` package is special:
1. Must have a `main()` function
2. This function is the entry point
3. Cannot be imported by other packages
4. Produces an executable binary

```go
package main  // This package will be an executable

import "fmt"

func main() {  // Entry point - required in main package
    fmt.Println("Hello, World!")
}
```

**Node.js comparison:**
```typescript
// There's no "special" main in Node.js
// The file you run with `node` is the entry point
console.log("Hello, World!")
```

---

## Adding and Removing Packages

### Adding a Dependency

**Node.js:**
```bash
npm install express
```

**Go:**
```bash
# Method 1: go get (downloads and adds to go.mod)
go get github.com/gin-gonic/gin

# Method 2: Just import it in code, then run:
go mod tidy  # Downloads missing deps, removes unused
```

### Removing a Dependency

**Node.js:**
```bash
npm uninstall express
```

**Go:**
```bash
# Just remove the import from your code, then:
go mod tidy  # Removes it from go.mod automatically
```

**Note:** `go mod tidy` is your best friend. Run it whenever you change dependencies!

---

## Exercises

1. Try changing `HelloHandler` to `helloHandler` in `handlers/hello.go` and see the compile error
2. Add a private helper function to `config.go` and try calling it from `main.go`
3. Run `go mod tidy` and examine the `go.sum` file
4. Try importing a package with an alias

---

## Next Steps

- Read **03-environment-variables.md** to see how Go handles configuration
- Try creating your own package in `pkg/` directory
