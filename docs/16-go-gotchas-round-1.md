# Go gotchas — Round 1 (basics)

A reference for the nuances that bit us while working through `quicky.go`. Each section has the rule, an example, and **why** Go is designed this way — because once you know the reason, the rule sticks.

> Complement to [`15-fmt-and-logging.md`](./15-fmt-and-logging.md), which covers `fmt` and `slog` separately.

---

## 1. Zero values & nil

Every variable in Go is initialized to its **zero value** — never garbage. The zero value depends on the type.

| Type | Zero value | Prints as | Safe to read? | Safe to write? |
|---|---|---|---|---|
| `string` | `""` | (empty) | ✅ | ✅ (but strings are immutable) |
| `int`, `float64`, `complex64`, … | `0`, `0.0` | `0` | ✅ | ✅ |
| `bool` | `false` | `false` | ✅ | ✅ |
| `*T` (pointer) | `nil` | `<nil>` | ❌ dereferencing panics | n/a |
| `[]T` (slice) | `nil` | `[]` | ✅ (len=0, range is a no-op, `append` works) | ❌ indexing panics |
| `map[K]V` | `nil` | `map[]` | ✅ (reads return zero value) | ❌ **assignment panics** |
| `chan T` | `nil` | `<nil>` | ❌ send/recv block **forever** | ❌ same |
| `interface{}` / `any` | `nil` | `<nil>` | n/a | n/a |
| `func` | `nil` | `<nil>` | ❌ calling panics | n/a |
| struct | each field at its zero value | `{...}` | ✅ | ✅ |

**Mental model:** `nil` only exists for *reference-ish* types (pointer, slice, map, channel, interface, function). Primitive types have concrete zero values like `0` and `false`.

### Why `nil` maps panic on write but nil slices don't

| | nil map | nil slice |
|---|---|---|
| Read | returns zero value silently | len=0, range yields nothing |
| Write | `m["x"] = 1` → **panic** | `append(s, 1)` returns a new slice (works) |

Slices are essentially `{ptr, len, cap}` — `append` knows how to allocate a backing array when needed. Maps have a much more complex internal structure (hash table with buckets) — you must call `make` to allocate it before writing.

---

## 2. Slices — length, capacity, and the `append` trap

```go
sl := make([]int, 3, 5)
// len=3, cap=5, slice points into a backing array of size 5
// indices 0,1,2 are zero-valued ints; indices 3,4 exist in cap but not in len
```

| `make` form | len | cap | When to use |
|---|---|---|---|
| `make([]int, 5)` | 5 | 5 | start with 5 zeroes you'll overwrite by index |
| `make([]int, 0, 5)` | 0 | 5 | empty slice you'll grow via `append`, preallocated to avoid realloc until you exceed 5 |
| `[]int{1, 2, 3}` | 3 | 3 | known elements at construction |
| `var s []int` | 0 | 0 | nil slice; valid, append starts allocating from scratch |

### The `append` gotcha

```go
func bad(s []int) {
    s = append(s, 99)  // local rebind only
}

s := make([]int, 0, 2)
bad(s)
fmt.Println(s) // []  — surprise!
```

**Why:** a slice is a small header `{ptr, len, cap}` passed *by value*. `append` may allocate a new backing array if `len == cap`, and even if it doesn't, the new `len` lives in the local header. The caller's header is unchanged.

**Fix:** either return the new slice, or take `*[]int`:
```go
func good(s []int) []int { return append(s, 99) }
s = good(s)
```

The first place this bites in production code: building up state in a helper.

---

## 3. Maps — silent zero-value lookups & rehashing

### Missing key returns the zero value, silently

```go
m := map[string]int{"alice": 30}
fmt.Println(m["bob"])  // prints 0, no error
```

This is the source of *so many* silent bugs. Use the **comma-ok** form when "key exists" matters:
```go
v, ok := m["bob"]
if !ok {
    // bob was missing — handle it
}
```

### Map values are not addressable

```go
users := map[string]User{"alice": {ID: 1}}
users["alice"].PointerMethod()  // compile error
&users["alice"]                  // compile error
```

**Why:** Go maps are hash tables that **rehash** on writes. Entries move in memory. If Go let you take the address of a value, that pointer would silently become invalid after the next insert:
```go
p := &m["alice"]   // hypothetical
m["bob"] = ...     // triggers rehash → "alice" now lives elsewhere
*p = ...           // writes to freed memory or to a different key!
```
Rather than expose this footgun, the language forbids it at compile time.

**Escape hatch:** store pointers in the map, not values.
```go
users := map[string]*User{"alice": {ID: 1}}
users["alice"].PointerMethod()  // works
users["alice"].Active = false   // mutates the User the pointer references
```
The `*User` pointers in the map can be moved on rehash (they're just words), but the `User` structs they point to live elsewhere and don't move.

---

## 4. Method receivers — value vs pointer

```go
type Counter struct{ n int }
func (c Counter)  IncVal() { c.n++ }   // value receiver — works on a copy
func (c *Counter) IncPtr() { c.n++ }   // pointer receiver — mutates the original
```

| Choose value receiver when | Choose pointer receiver when |
|---|---|
| The method doesn't mutate state | The method **mutates** the receiver |
| The type is small (≤ ~8 bytes, no slice/map fields) | The type is "big-ish" (multiple fields, embedded structs) |
| The type is logically immutable (e.g. `time.Time`) | *Any* other method on the type already uses pointer receiver (be consistent) |

**Rule of thumb when in doubt: pointer.** Mixing value and pointer receivers on the same type is confusing and sometimes outright broken — see "method sets" below.

### Method sets — the subtle rule

| Receiver type | Method set of `T` includes | Method set of `*T` includes |
|---|---|---|
| `func (t T) ...` | ✅ | ✅ |
| `func (t *T) ...` | ❌ | ✅ |

In plain English: **a pointer can call any method on T or *T. A value can only call methods defined on T.** That's why this fails:

```go
type S struct{}
func (s *S) Method() {}

var x S
x.Method()              // works — Go auto-takes &x because x is addressable

m := map[string]S{"k": {}}
m["k"].Method()         // compile error — map values aren't addressable
                        // Go can't auto-take the address
```

---

## 5. Addressability — the unified rule

Several language rules trace back to one concept: **does this expression have a stable memory address the compiler can use?**

| Expression | Addressable? | Why |
|---|---|---|
| `x` (local variable) | ✅ | named storage in a stack frame |
| `*p` (dereferenced pointer) | ✅ | `p` points to addressable memory |
| `s[i]` (slice element) | ✅ | backing array doesn't relocate on its own |
| `arr[i]` (array element, when array is addressable) | ✅ | same |
| `someStruct.field` (when `someStruct` is addressable) | ✅ | field has a stable offset |
| `m["k"]` (map value) | ❌ | rehashing relocates entries |
| `User{}` (composite literal) | ❌ | unnamed temporary, no slot allocated |
| `f()` (function return value) | ❌ | temporary until you assign it |
| `42` (constant) | ❌ | constants are values, not memory |
| `"hello"` (string literal) | ❌ | same |

**Consequences:**
1. You can only do `&expr` if `expr` is addressable.
2. You can only call a pointer-receiver method via `x.Method()` shorthand if `x` is addressable. Otherwise: assign to a variable first, or use `(&Literal{}).Method()` explicitly.
3. You can only assign to an expression if it's addressable.

This *one rule* explains: nil map writes (different story — see §1), map-value `&` errors, "cannot take the address of `User{}`", and why `f().field = x` doesn't compile.

---

## 6. The `String() string` method — the `Stringer` interface

```go
type Stringer interface { String() string }
```

Anything with a `String() string` method satisfies `fmt.Stringer`. `fmt.Println` (and friends) check for it and call it. **This is why your `User(2, Cena2)` output worked.**

### Receiver choice for `String()`

| Receiver | `fmt.Println(u)` where `u` is `User` | `fmt.Println(&u)` where `&u` is `*User` |
|---|---|---|
| `func (u User) String() string` | ✅ called | ✅ called (via auto-deref) |
| `func (u *User) String() string` | ❌ NOT called (prints raw struct) | ✅ called |

**Convention:** use value receiver for `String()` unless the method mutates (it almost never should). Same applies to `MarshalJSON`, `Error()`, etc.

---

## 7. Errors — wrapping and inspecting

### Wrapping with `%w`

| Form | What it does |
|---|---|
| `return err` | propagate unchanged |
| `fmt.Errorf("ctx: %w", err)` | **wrap** — preserves error chain, adds context |
| `fmt.Errorf("ctx: %v", err)` | flatten to string — chain is lost |
| `fmt.Errorf("ctx: %s", err)` | same as `%v` for errors — chain lost |

**Default to `%w` unless you have a reason to hide the cause.**

### Inspecting wrapped errors

| Function | Asks | Use for |
|---|---|---|
| `errors.Is(err, target)` | "Is this **specific** sentinel error anywhere in the chain?" | `io.EOF`, `sql.ErrNoRows`, `context.Canceled`, `context.DeadlineExceeded` |
| `errors.As(err, &target)` | "Is there an error of this **type** in the chain? Bind it to `target`." | typed errors with fields (like `*apperror.AppError`) |

Real example from this project (`internal/handler/user.go`):
```go
var appErr *apperror.AppError
if errors.As(err, &appErr) {
    response.WriteAppError(w, appErr)  // appErr.Code/Status/Fields all available
    return
}
```

---

## 8. The `nil` interface trap (preview — bites later)

```go
var err error
var typedErr *MyError = nil
err = typedErr
fmt.Println(err == nil)   // false  — surprise!
```

**Why:** an interface holds two things — a *type* and a *value*. `err` got the type `*MyError` (even though the value is nil), so `err == nil` is false. You'll meet this when functions return typed nil errors.

**Fix:** return `nil` explicitly when there's no error.
```go
func work() error {
    var e *MyError
    if /* something went wrong */ {
        e = &MyError{...}
        return e
    }
    return nil  // ← not `return e`
}
```

(We'll cover this properly in Round 2 when we talk about interfaces and typed errors.)

---

## TL;DR — the rules to memorize before Round 2

1. **`nil` map writes panic, `nil` slice writes don't** — `make` your maps.
2. **Map lookups for missing keys return zero values silently** — use comma-ok when you need to know.
3. **`append` may return a new backing array** — always assign the result back: `s = append(s, x)`.
4. **Map values are not addressable** (because of rehashing). Use `map[K]*V` when you need pointer access.
5. **Pointer-receiver method calls need an addressable value.** Literals and map values aren't addressable.
6. **`String()` (and `Error()`, `MarshalJSON`, etc.) → value receiver unless mutating.**
7. **Wrap errors with `%w`, inspect with `errors.Is` / `errors.As`.**
