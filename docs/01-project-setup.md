# 01 - Project Setup

## Go Modules: The Package Manager

### What is go.mod?

In Node.js, you have `package.json` to track your project and dependencies. In Go, you have `go.mod`.

**Node.js:**
```json
{
  "name": "task-manager",
  "version": "1.0.0",
  "dependencies": {
    "express": "^4.18.0"
  }
}
```

**Go (go.mod):**
```go
module github.com/user/task-manager

go 1.21

require github.com/joho/godotenv v1.5.1
```

### Creating a New Module

**Node.js:**
```bash
npm init
# or
npm init -y
```

**Go:**
```bash
go mod init github.com/user/task-manager
```

The module path (`github.com/user/task-manager`) is important:
- It's how other Go code imports your packages
- It typically matches your repository URL
- For local/learning projects, you can use any path

### go.sum vs package-lock.json

Both track exact versions of dependencies for reproducible builds:

| File | Purpose |
|------|---------|
| `package-lock.json` | Tracks npm package versions |
| `go.sum` | Tracks Go module checksums (more secure!) |

Go's approach is more secure because it stores cryptographic checksums, not just versions.

---

## Project Structure

### Why This Structure?

Go doesn't enforce a project structure, but there are conventions from the community:

```
task-manager/
├── cmd/           # Application entry points
├── internal/      # Private packages (can't be imported externally)
├── pkg/           # Public packages (can be imported)
└── docs/          # Documentation
```

### cmd/ Directory

Contains your main applications. Each subdirectory becomes a separate binary:

```
cmd/
├── api/           # HTTP API server -> builds to 'api' binary
├── cli/           # CLI tool -> builds to 'cli' binary
└── worker/        # Background worker -> builds to 'worker' binary
```

**Why?** One repository can produce multiple programs. NestJS does this differently with monorepos and workspaces.

### internal/ Directory

This is SPECIAL in Go! Any package inside `internal/` cannot be imported by code outside your module.

```go
// This works (inside your module):
import "github.com/user/task-manager/internal/handlers"

// This FAILS (from another module):
import "github.com/user/task-manager/internal/handlers"
// Error: use of internal package not allowed
```

**Why?** It's Go's way of creating truly private code. In Node.js, everything in `node_modules` is technically accessible.

### pkg/ Directory

Packages meant to be imported by other projects. If you're building a library, put reusable code here.

---

## Building and Running

### Development

**Node.js:**
```bash
# Run directly (interpreted)
node index.js

# Or with TypeScript
npx ts-node src/index.ts

# Or compile first
tsc && node dist/index.js
```

**Go:**
```bash
# Run directly (compiles in temp dir, then runs)
go run ./cmd/api

# Or build first, then run
go build -o bin/task-manager ./cmd/api
./bin/task-manager
```

### Key Difference: Compilation

| Language | Type | Result |
|----------|------|--------|
| JavaScript | Interpreted | Runs in Node.js runtime |
| TypeScript | Transpiled | Converts to JS, runs in Node.js |
| Go | Compiled | Single binary, no runtime needed |

A Go binary includes everything needed to run. No `node_modules`, no runtime installation on the server. Just copy the binary and run it!

---

## Commands Reference

| Task | Node.js | Go |
|------|---------|-----|
| Initialize project | `npm init` | `go mod init <path>` |
| Add dependency | `npm install <pkg>` | `go get <pkg>` |
| Install all deps | `npm install` | `go mod download` |
| Remove unused deps | `npm prune` | `go mod tidy` |
| Update deps | `npm update` | `go get -u ./...` |
| Run | `npm start` / `node index.js` | `go run ./cmd/api` |
| Build | `npm run build` | `go build` |

---

## Next Steps

- Read **02-packages-and-modules.md** to understand Go's import system
- Try running `go mod tidy` to see it download dependencies
- Explore the generated `go.sum` file
