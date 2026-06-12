# `fmt` & logging — quick reference

## 1. The `fmt` family: which function to call

| Function | Returns | Adds newline? | Goes to | Use when |
|---|---|---|---|---|
| `fmt.Print(a, b)` | `(n int, err error)` | ❌ | stdout | rarely — you almost always want `Println` or `Printf` |
| `fmt.Println(a, b)` | `(n int, err error)` | ✅ auto | stdout | quick debugging, simple prints; args separated by spaces |
| `fmt.Printf(fmt, args)` | `(n int, err error)` | ❌ — add `\n` yourself | stdout | formatted output; full control over layout |
| `fmt.Sprint(a, b)` | `string` | ❌ | nothing printed | build a string by concatenating args |
| `fmt.Sprintln(a, b)` | `string` | ✅ (in the returned string) | nothing printed | rarely useful |
| `fmt.Sprintf(fmt, args)` | `string` | ❌ | nothing printed | **build a formatted string** — used everywhere: log messages, error messages, building SQL/URLs (carefully) |
| `fmt.Errorf(fmt, args)` | `error` | n/a | nothing printed | build an `error` value, often with `%w` to wrap an existing error |
| `fmt.Fprint(w, ...)` | `(n int, err error)` | ❌ | the `io.Writer` you pass | write to a file, an `http.ResponseWriter`, a buffer |
| `fmt.Fprintln(w, ...)` | `(n int, err error)` | ✅ | the `io.Writer` you pass | same as above with newline |
| `fmt.Fprintf(w, fmt, args)` | `(n int, err error)` | ❌ | the `io.Writer` you pass | formatted write to any `io.Writer` |

**Mental model:**
- `Print*` → stdout
- `Sprint*` → returns string, prints nothing
- `Fprint*` → writes to any `io.Writer` (file, response, buffer)
- `Errorf` → returns an `error`
- Suffix `f` = format string with verbs (`%d`, `%s`, …)
- Suffix `ln` = adds newline automatically and separates args with spaces

## 2. Format verbs — what each `%` does

| Verb | Meaning | Example input | Example output |
|---|---|---|---|
| `%v` | default format (works on anything) | `User{1, "alice"}` | `{1 alice}` |
| `%+v` | default format **with field names** | `User{1, "alice"}` | `{ID:1 Username:alice}` |
| `%#v` | Go-syntax representation (for debugging) | `User{1, "alice"}` | `main.User{ID:1, Username:"alice"}` |
| `%T` | the **type** of the value | `42` | `int` |
| `%d` | integer, base 10 | `42` | `42` |
| `%b` | integer, base 2 | `5` | `101` |
| `%o` | integer, base 8 | `8` | `10` |
| `%x` / `%X` | integer, base 16 (lower/upper) | `255` | `ff` / `FF` |
| `%c` | character (rune) | `65` | `A` |
| `%s` | string (or anything implementing `String()`) | `"hi"` | `hi` |
| `%q` | quoted string | `"hi"` | `"hi"` |
| `%f` | float, decimal | `3.14` | `3.140000` |
| `%.2f` | float with 2 decimals | `3.14159` | `3.14` |
| `%e` | scientific notation | `1234.5` | `1.234500e+03` |
| `%g` | "smart" float — shortest accurate | `0.0001` | `0.0001` |
| `%t` | boolean | `true` | `true` |
| `%p` | pointer address | `&x` | `0xc0000140a0` |
| `%w` | **wrap an error** (only valid in `fmt.Errorf`) | `fmt.Errorf("x: %w", err)` | builds a chained error |
| `%%` | literal percent sign | — | `%` |

**Tip:** when in doubt, use `%v`. For errors in `fmt.Errorf`, always use `%w` (not `%v`/`%s`) so `errors.Is` / `errors.As` can walk the chain.

## 3. `%w` vs `%v` vs `%s` for errors

| Form | What it does | Use when |
|---|---|---|
| `fmt.Errorf("ctx: %w", err)` | **Wraps**: preserves the chain | almost always, when propagating an error with extra context |
| `fmt.Errorf("ctx: %v", err)` | Formats `err` as a string into a new error — **chain is lost** | rarely — only if you intentionally want to hide the underlying error |
| `fmt.Errorf("ctx: %s", err)` | same as `%v` for errors — chain is lost | almost never; treat as a smell |
| `return err` | Pass through unchanged | when you have no extra context to add |

## 4. `errors.Is` vs `errors.As`

| Function | Question it answers | Use for |
|---|---|---|
| `errors.Is(err, target)` | "Is this *exact* sentinel error somewhere in the chain?" | sentinel values: `io.EOF`, `sql.ErrNoRows`, `context.Canceled`, `context.DeadlineExceeded` |
| `errors.As(err, &target)` | "Is there an error of this *type* in the chain? If yes, bind it." | typed errors with fields (like `*apperror.AppError`) |

Example from this project:
```go
var appErr *apperror.AppError
if errors.As(err, &appErr) {
    // appErr.Code, appErr.Status, appErr.Fields all accessible here
}
```

## 5. `fmt` vs `slog` — when to use which

| Scenario | Use | Why |
|---|---|---|
| Quick debug while iterating | `fmt.Println` | shortest, throwaway |
| Building a string to use elsewhere | `fmt.Sprintf` | returns a value, doesn't print |
| Building an error | `fmt.Errorf` with `%w` | preserves the error chain |
| Application logging in production code | `slog` (`slog.Info`, `slog.Error`) | structured, JSON-able, level-controlled, attaches context |
| Logging an HTTP request | `slog.Info("request", "method", m, "path", p, ...)` | structured fields, machine-parseable |
| Logging a panic / fatal startup error | `slog.Error(...)` + `os.Exit(1)` | structured; never use `log.Fatal` in real services |

**Rules of thumb:**
- In production code, use `slog`, not `fmt.Println`.
- `slog` arguments come in **key/value pairs**: `slog.Info("msg", "key1", val1, "key2", val2)`.
- Never concatenate values into the message: ❌ `slog.Info("port " + port)`. Use fields: ✅ `slog.Info("server starting", "port", port)`. The whole point of `slog` is structured output you can query later.

## 6. `slog` cheat sheet

| Call | When |
|---|---|
| `slog.Debug("msg", "k", v)` | verbose dev-only info (only printed if level is Debug) |
| `slog.Info("msg", "k", v)` | normal operational events ("server started", "request handled") |
| `slog.Warn("msg", "k", v)` | recoverable problems ("retrying after timeout") |
| `slog.Error("msg", "error", err, "k", v)` | unexpected failures; pair with returning/handling the error |
| `slog.With("requestID", id).Info(...)` | bind common fields to a logger so every line includes them |
| `slog.SetDefault(logger)` | set the package-level default logger (done once at startup) |

**Project example:** see `cmd/api/main.go` — `setLogger()` builds a JSON handler and sets it as default. Every `slog.X` call in the codebase then writes JSON to stdout, with level controlled by env.

## 7. Common newbie traps

| Trap | Fix |
|---|---|
| `fmt.Printf("hello %s", name)` produces output on the same line as the next print | add `\n`: `fmt.Printf("hello %s\n", name)` |
| `fmt.Println("port " + strconv.Itoa(p))` | use `fmt.Println("port", p)` (Println handles spacing) |
| `fmt.Errorf("bad: %s", err)` | use `%w`: `fmt.Errorf("bad: %w", err)` |
| `slog.Info("port " + port)` | structured: `slog.Info("server starting", "port", port)` |
| Forgetting `%` count doesn't match args | `go vet` catches this — run it |
| Using `log.Fatal` / `log.Panic` in a server | almost always wrong — use `slog.Error` + `os.Exit(1)` in `main` only |
