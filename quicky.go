package main

import (
	"encoding/json"
	"errors"
	"fmt"
)


type User struct {
	ID int
	Username string
	Active bool
}

type Users []User

func (users Users) activeUsers() Users {
	activeUsers := make([]User, 0, len(users));
	for i := range users {
		if users[i].Active {
			activeUsers = append(activeUsers, users[i])
		}
	}
	return activeUsers;
}


func(u User) String() string {
	return fmt.Sprintf("User(%d, %s)", u.ID, u.Username)
}
// ============================================================================
// ROUND 1: BASICS
// Run with: go run quicky.go
// Tip: comment out calls in main() you haven't done yet.
// ============================================================================

func main() {
	ex1ZeroValues()
	ex2Slices()
	ex3Maps()
	ex4Structs()
	ex5FmtFamily()
	ex6Functions()
	ex7Pointers()
	// ROUND 2
	ex8Defer()
	ex9PanicRecover()
	ex10CustomError()
	ex11Interfaces()
	ex12Enums()
}

// ----------------------------------------------------------------------------
// EX 1 — Zero values
// Declare ONE variable of each type below (use `var`, not `:=`) and print them
// all on one line with fmt.Println. Then try to write to the nil map — observe
// the panic, then fix it by initializing the map with make() before writing.
//
// Types to declare: string, int, bool, float64, *int, []int, map[string]int
//
// Expected output (before the panic-fix):
//   "" 0 false 0 <nil> [] map[]
// ----------------------------------------------------------------------------
func ex1ZeroValues() {
	fmt.Println("--- ex1ZeroValues ---")
	// TODO
	var s string
	var ss int
	var fl float64
	var booll bool
	var p *int
	var pp []int

	// var mappx map[string]int
	// mappx["x"] = 1

	var mapp = make(map[string]int)

	fmt.Println(s, ss, fl,booll, p, pp, mapp);

	mapp["alice"] = 1;

	mapp["x"] = 1;

	fmt.Println(mapp);
}

// ----------------------------------------------------------------------------
// EX 2 — Slices
// 1. Create a slice of ints: [10, 20, 30]
// 2. Append 40 to it.
// 3. Iterate with `for i, v := range slice` and print "index=i value=v".
// 4. Compute and print the sum.
// 5. Bonus: create an empty slice with make([]int, 0, 5) — what's its len? cap?
//    Print them: fmt.Println(len(s), cap(s))
// ----------------------------------------------------------------------------
func ex2Slices() {
	fmt.Println("--- ex2Slices ---")
	slice1 := []int{10,20,30}
	slice1 = append(slice1, 40)

	sum := 0;
	for key,value := range slice1 {
		fmt.Println(key, "---", value)
		sum += value
	}
	fmt.Println(sum);

	slice2 := make([]int, 5);
	fmt.Println(len(slice2), cap(slice2))

	slice3 := make([]int, 0, 5);
	fmt.Println(len(slice3), cap(slice3))
}

// ----------------------------------------------------------------------------
// EX 3 — Maps
// 1. Create map[string]int{"alice": 30, "bob": 25}.
// 2. Add "carol": 28.
// 3. Look up "alice" — print her age.
// 4. Use the "comma-ok" form to check if "dave" exists; print "not found" if not.
// 5. Delete "bob".
// 6. Range over the map and print each "name=age".
//
// Reminder: map iteration order is NOT guaranteed. Don't rely on it.
// ----------------------------------------------------------------------------
func ex3Maps() {
	fmt.Println("--- ex3Maps ---")
	ages := map[string]int{"alice": 30, "bob": 25}
	ages["carol"] = 28;

  	fmt.Println("alice:", ages["alice"])

	
	if _, ok := ages["dave"]; !ok {
		fmt.Println("not found")
	}

	delete (ages, "bob");
	for k,v := range ages {
		fmt.Println("Name", k, "age", v);
	}
}

// ----------------------------------------------------------------------------
// EX 4 — Structs
// 1. Define a struct type User with fields: ID int, Username string, Active bool.
//    (Define it OUTSIDE this function, at package level — see the TODO below.)
// 2. Create a []User with 3 users (mix of active/inactive).
// 3. Write a function activeUsers(users []User) []User that returns only the
//    active ones. Call it and print the result.
// 4. Bonus: add a method (u User) String() string that returns
//    fmt.Sprintf("User(%d, %s)", u.ID, u.Username). Now Println(u) uses it.
// ----------------------------------------------------------------------------

// TODO: define `type User struct { ... }` here at package level
// TODO: define `func activeUsers(users []User) []User { ... }` here

func ex4Structs() {
	fmt.Println("--- ex4Structs ---")
	users := Users{
		{ID: 1, Username: "Cena1", Active: false },
		{ID: 2, Username: "Cena2", Active: true },
		{ID: 3, Username: "Cena3", Active: true },
		{ID: 4, Username: "Cena4", Active: false },
	}
	activeUsers := users.activeUsers()
	fmt.Println(activeUsers);

	fmt.Println(users[1])
}

// ----------------------------------------------------------------------------
// EX 5 — The fmt family
// You constantly mix these up. Make them concrete:
//
// 1. fmt.Println("a", "b")          — prints, adds spaces + newline → STDOUT
// 2. fmt.Printf("hi %s\n", "you")   — formatted, you control everything → STDOUT
// 3. s := fmt.Sprintf("hi %s", "you") — RETURNS a string, prints NOTHING
// 4. fmt.Errorf("bad: %w", err)     — RETURNS an error, used for wrapping
//
// Tasks:
// a. Use Sprintf to build a string "user 42 has 7 messages" and print it.
// b. Use Printf with %d, %s, %v to print: 42, "alice", true on one line.
// c. Use Printf with %T to print the TYPE of a slice, a map, and a *int.
// ----------------------------------------------------------------------------
func ex5FmtFamily() {
	fmt.Println("--- ex5FmtFamily ---")
	fmt.Println("a", "b")
	fmt.Printf("hi %s\n", "you")
	s := fmt.Sprintf("hi %s\n", "you")
	fmt.Println(s)
	err := fmt.Errorf("bad: %s", "errorrr")
	fmt.Println(err)

	// tasks
	ss := fmt.Sprintf("user %d has %d messages", 42, 7)
	fmt.Println(ss);

	fmt.Printf("%d %s %v \n", 42, "alice", true, )

	var x int
	sl := []int{}
	mapp := make(map[string]int)

	fmt.Printf("%T %T %T\n", sl, mapp, &x)
}

// ----------------------------------------------------------------------------
// EX 6 — Functions: multiple returns
// 1. Write divide(a, b float64) (float64, error) that returns (a/b, nil) when
//    b != 0, and (0, fmt.Errorf("divide by zero")) when b == 0.
// 2. Call it twice: once with valid args, once with b=0. Handle the error idiomatically:
//        result, err := divide(...)
//        if err != nil { fmt.Println("error:", err); ... }
// 3. Bonus: write minmax(nums []int) (min, max int) using NAMED return values
//    and a `return` (naked) at the end.
// ----------------------------------------------------------------------------

// TODO: define divide here
// TODO: define minmax here
func divide(a float64, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("divide by zero");
	}
	return (a / b), nil;
}

func minmax(nums []int) (min, max int) {
	if(len(nums) == 0) {
		return
	}
	min , max = nums[0], nums[0]
	for _, v := range nums[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return;
}

func ex6Functions() {
	fmt.Println("--- ex6Functions ---")
	fmt.Println(divide(10, 2));
	fmt.Println(divide(10,0));
	
	fmt.Println(minmax([]int{1,2,3,4,5,6,7}));
}

// ----------------------------------------------------------------------------
// EX 7 — Pointers
// 1. Write a function double(n int) that takes int and doubles it. Call it with
//    x := 5; double(x); print x. Observe that x is UNCHANGED.
// 2. Write doublePtr(n *int) that doubles via pointer. Call doublePtr(&x).
//    Print x — now changed.
// 3. Write a function reset(u *User) that sets u.Active = false. Call it on
//    a User and print before/after.
// 4. Bonus question (no code): why does `(&User{}).String()` work but
//    `User{}.PointerOnlyMethod()` might not? (Hint: addressability.)
// ----------------------------------------------------------------------------

// TODO: define double, doublePtr, reset

func double(x int) int {
	return x * 2
}

func doublePointer(x *int) {
	*x = *x * 2
}

func reset(u *User) *User {
	u.Active = false
	return u
}

// ============================================================================
// ROUND 2 starts after ex7Pointers
// ============================================================================

func ex7Pointers() {
	fmt.Println("--- ex7Pointers ---")
	x := 5
	doubleIt := double(x)
	fmt.Println(doubleIt, x)

	doublePointer(&x)
	fmt.Println(x)

	user := &User{
		Active: true,
		Username: "Blah balh",
		ID: 1,
	}

	fmt.Println(user.Active)
	reset(user);
	fmt.Println(user.Active)
}

// ============================================================================
// ROUND 2: defer, panic/recover, custom errors, interfaces
// 4 exercises. Same drill — uncomment in main(), implement, run.
// ============================================================================

// ----------------------------------------------------------------------------
// EX 8 — defer (the most useful Go keyword you'll use daily)
//
// defer schedules a function call to run when the SURROUNDING FUNCTION returns
// (not when the enclosing block ends — defer is function-scoped).
//
// Three things you MUST internalize:
//
//   (a) ORDER: defers run LIFO (last deferred → first to run).
//       Try:
//           defer fmt.Println("a")
//           defer fmt.Println("b")
//           defer fmt.Println("c")
//       Predict the output BEFORE running. Then run.
//
//   (b) ARGUMENTS ARE EVALUATED AT DEFER TIME, not at run time.
//       Try:
//           x := 1
//           defer fmt.Println("deferred x =", x)
//           x = 99
//       What prints? The captured x is 1, NOT 99 — fmt.Println's args were
//       evaluated when `defer` was executed.
//
//   (c) defer inside a LOOP accumulates — they all fire when the function
//       returns, NOT at the end of each iteration. This is a common leak:
//           for _, f := range files {
//               h, _ := os.Open(f)
//               defer h.Close()  // ← every Close stacks up; held open until function exits
//               // ... use h ...
//           }
//       Tasks:
//   1. Build the (a) example — 3 defers, predict, then verify.
//   2. Build the (b) example with x:=1; print BEFORE-and-AFTER style.
//   3. Bonus: write a tiny function that returns a NAMED int and uses
//      `defer func() { ret *= 2 }()` to double it. Verify the deferred
//      mutation reaches the caller.
// ----------------------------------------------------------------------------

func namedInt() (ret int) {
	defer fmt.Println("by arg:", ret)         // argument ret evaluated NOW → 0                                                                                                             
      defer func() {                                                                                                                                                                          
          fmt.Println("by closure:", ret)       // body reads ret at exit → 5
		  ret *= 2                                                                                                            
      }()                                                                                                                                                                                     
      return 5 
}

func ex8Defer() {
	fmt.Println("--- ex8Defer ---")
	defer fmt.Println("a")
	defer fmt.Println("b")
	defer fmt.Println("c")

	x := 1
	defer fmt.Println(x)
	x = 99

	i := 0;
	for i <= 3 {
		defer fmt.Println("loop", i)
		i++
	}

	fmt.Println(namedInt())
}

// ----------------------------------------------------------------------------
// EX 9 — panic / recover (use sparingly)
//
// panic stops normal flow and starts unwinding (running deferred calls).
// recover catches a panic — but ONLY when called from a DEFERRED function.
//
// In real Go code, panic is used:
//   - for programmer-error invariants ("this should never happen")
//   - in libraries that recover at their public boundary (e.g. chi's Recoverer)
// You rarely panic yourself in app code. You return errors instead.
//
// Tasks:
//   1. Write a function safeCall(fn func()) that runs fn() and recovers from
//      any panic, printing "recovered:" + the panic value. Signature:
//          func safeCall(fn func()) { ... }
//      Inside, use:
//          defer func() {
//              if r := recover(); r != nil {
//                  fmt.Println("recovered:", r)
//              }
//          }()
//          fn()
//   2. Call safeCall with a function that does `panic("boom")` — observe.
//   3. Call safeCall with a function that does NOT panic — observe (no print).
//
// Question to answer in your head: WHY must recover() be inside a deferred call?
// (Hint: the panic is mid-unwind. Only deferreds run during unwinding.)
// ----------------------------------------------------------------------------

func safeCall(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recoverd error: \n", r)
		}
	}()

	fn()
}
func ex9PanicRecover() {
	fmt.Println("--- ex9PanicRecover ---")
	// TODO
	fn := func() {
		fmt.Println("I'll do panic")
		panic("boom panic")
	}
	
	safeCall(fn)

	safefn := func() {
		fmt.Println("I'll do safe")
	}
	safeCall(safefn)
}


// ----------------------------------------------------------------------------
// EX 10 — Custom error type + errors.As + the nil interface trap
//
// Errors in Go are just values that satisfy the `error` interface:
//     type error interface { Error() string }
//
// Tasks:
//   1. Define a type NotFoundError with field Resource string. Implement
//      Error() string returning fmt.Sprintf("%s not found", e.Resource).
//      (Use a POINTER receiver: `func (e *NotFoundError) Error() string`.)
//
//   2. Write a function findUser(id int) (*User, error). If id == 0,
//      return (nil, &NotFoundError{Resource: "user"}). Otherwise return
//      (&User{ID: id, Username: "fake"}, nil).
//
//   3. Call findUser(0). Check the error with errors.As:
//          var nfe *NotFoundError
//          if errors.As(err, &nfe) {
//              fmt.Println("got NotFound for:", nfe.Resource)
//          }
//      Import "errors" at the top.
//
//   4. The nil-interface trap. Write a function badNoError() error that does:
//          var e *NotFoundError = nil
//          return e
//      Then in main:
//          err := badNoError()
//          fmt.Println(err == nil)  // prints FALSE — surprise!
//      Why? An interface holds (type, value). Even though the *value* is nil,
//      the *type* is *NotFoundError, so the interface isn't nil.
//      The fix is to `return nil` explicitly when there's no error.
// ----------------------------------------------------------------------------

// TODO: define type NotFoundError + Error() method
// TODO: define findUser
// TODO: define badNoError (for the nil-interface trap)

type NotFoundError struct {
	Resource string
}

func(e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Resource)
}

func findUser(id int) (*User, error) {
	if(id == 0) {
		return nil, &NotFoundError{Resource:  string(id)}
	}

	return &User{ID: id, Username: "hahhaa", Active: true}, nil
}

func badNoError() error {
	var e *NotFoundError = nil
	return e
}

func ex10CustomError() {
	fmt.Println("--- ex10CustomError ---")
	var nfe *NotFoundError
	user,err := findUser(0)
	if errors.As(err, &nfe) {
		fmt.Println("got NotFound for:", nfe.Resource)
	}
	fmt.Println(user)

	errr := badNoError()
	fmt.Println(errr == nil)
}

// ----------------------------------------------------------------------------
// EX 11 — Interfaces (the real heart of Go)
//
// An interface is a SET OF METHODS. Any type with those methods satisfies it
// — no `implements` keyword. This is "structural / implicit" satisfaction.
//
// Tasks:
//   1. Define an interface Greeter:
//          type Greeter interface { Greet() string }
//
//   2. Make TWO types satisfy it:
//          type English struct{ Name string }
//          func (e English) Greet() string { return "Hello, " + e.Name }
//
//          type Spanish struct{ Name string }
//          func (s Spanish) Greet() string { return "Hola, " + s.Name }
//
//   3. Write a function sayAll(gs []Greeter) that iterates and prints each
//      Greet(). Call it with a mixed slice:
//          sayAll([]Greeter{ English{"Alice"}, Spanish{"Bob"} })
//
//   4. Type assertion: given a Greeter g, try to recover the concrete English:
//          if e, ok := g.(English); ok {
//              fmt.Println("it was English named", e.Name)
//          }
//
//   5. Type switch: write a function describe(g Greeter) that prints
//      what concrete type it is:
//          switch v := g.(type) {
//          case English: fmt.Println("English:", v.Name)
//          case Spanish: fmt.Println("Spanish:", v.Name)
//          default:      fmt.Println("unknown")
//          }
//
//   6. Bonus — constants with iota (you'll use this for roles, statuses, etc.):
//          type Role int
//          const (
//              RoleAdmin Role = iota // 0
//              RoleUser              // 1
//              RoleGuest             // 2
//          )
//      Print all three. Add a String() method on Role so fmt.Println prints
//      "admin" / "user" / "guest" instead of 0/1/2.
// ----------------------------------------------------------------------------

// TODO: define Greeter, English, Spanish
// TODO: define sayAll, describe
// TODO: define Role + iota + String()

type Greeter interface {
	Greet() string
}

type English struct{ Name string }
func (e English) Greet() string { return "Hello, " + e.Name }

type Spanish struct{ Name string }
func (s Spanish) Greet() string { return "Hola, " + s.Name }


func sayAll(greeters []Greeter) {
	for _, v := range greeters {
		fmt.Println(v.Greet())
	}
}

func describe(g Greeter) {
	switch v := g.(type) {                                
      case English:
          fmt.Println("English greeter:", v.Name)
      case Spanish:                                                                                                                                                          
          fmt.Println("Spanish greeter:", v.Name)
      default:                                                                                                                                                               
          fmt.Println("unknown greeter:", g.Greet())        
      }                
}

func ex11Interfaces() {
	fmt.Println("--- ex11Interfaces ---")
	greeters := []Greeter{English{Name: "English"}, Spanish{Name: "Spanish"}}
	sayAll(greeters)

	var oneGreeter Greeter = English{Name:"English"}
	if e, ok := oneGreeter.(English); ok {
 		fmt.Println("got English:", e.Name)
	}

	if _, ok := oneGreeter.(Spanish); !ok {
      fmt.Println("not Spanish — correct")                                                                                                                                   
  	}
	
	describe(English{Name: "Alice"})
  	describe(Spanish{Name: "Bob"}) 
}

// ----------------------------------------------------------------------------
// EX 12 — "Enums" the Go way (vs TypeScript)
//
// Go has NO `enum` keyword. The idiomatic pattern is:
//   1. A custom integer (or string) type
//   2. Constants of that type, often using `iota`
//   3. (Optional) a String() method for readable output
//   4. (Optional) MarshalJSON / UnmarshalJSON for the wire format
//
// QUICK COMPARISON with TypeScript:
//
//   TypeScript:
//     enum Role { Admin, User, Guest }      // 0, 1, 2 + reverse mapping at runtime
//     enum Role { Admin = "admin", ... }    // string-valued; no reverse mapping
//     // The enum NAME is also a runtime object you can use.
//
//   Go:
//     type Role int
//     const (
//         RoleAdmin Role = iota
//         RoleUser
//         RoleGuest
//     )
//     // No runtime "Role" object. No reverse mapping. Just typed constants.
//     // For human-readable output, you add a String() method yourself.
//
// Tasks (in your file, not in your head):
//
//   1. Numeric enum with iota.
//      Define at PACKAGE LEVEL:
//          type Role int
//          const (
//              RoleAdmin Role = iota   // 0
//              RoleUser                // 1
//              RoleGuest               // 2
//          )
//      In ex12Enums(), print all three. They print as 0, 1, 2.
//
//   2. String() method so Println shows names, not numbers.
//      Add (also at package level):
//          func (r Role) String() string {
//              switch r {
//              case RoleAdmin: return "admin"
//              case RoleUser:  return "user"
//              case RoleGuest: return "guest"
//              default:        return "unknown"
//              }
//          }
//      Re-run ex12Enums and watch the same Println calls now print
//      "admin" / "user" / "guest" instead of 0 / 1 / 2.
//
//   3. Bitmask flags with `1 << iota` (a real Go idiom):
//          type Permission int
//          const (
//              PermRead Permission = 1 << iota  // 1
//              PermWrite                        // 2
//              PermDelete                       // 4
//              PermAdmin                        // 8
//          )
//      In ex12Enums:
//          p := PermRead | PermWrite
//          fmt.Println(p)                       // 3
//          fmt.Println(p & PermWrite != 0)      // true — bit is set
//          fmt.Println(p & PermDelete != 0)     // false
//
//   4. String-backed enum (often BETTER for things that hit JSON/DB):
//          type Status string
//          const (
//              StatusActive   Status = "active"
//              StatusInactive Status = "inactive"
//          )
//      Print StatusActive — no String() method needed; it's already a string.
//      Advantage: JSON marshal produces "active", not 0.
//
//   5. Bonus thought (no code): why prefer string-backed enums for anything
//      that crosses a boundary (JSON, DB, gRPC)?
//      Hint: if you insert a new value in the MIDDLE of an iota block later,
//      everything below it shifts by 1. Old stored 1s suddenly mean
//      something different. With string values, "active" stays "active"
//      forever, regardless of order changes.
// ----------------------------------------------------------------------------

// TODO: define Role + iota + String()
// TODO: define Permission + bitmask iota
// TODO: define Status + string constants
type Role int

const (
	RoleAdmin Role = iota
	RoleUser
	RoleGuest
)

type Status string 
const (
	StatusActive Status = "Active"
	StatusInactive Status = "inactive"
)

func(r Role) String() string {
	 switch r {
		case RoleAdmin: return "admin"
		case RoleUser:  return "user"
		case RoleGuest: return "guest"
		default:        return "unknown"
	}
}

 type Permission int
	const (
		PermRead Permission = 1 << iota  // 1
		PermWrite                        // 2
		PermDelete                       // 4
		PermAdmin                        // 8
	)

func ex12Enums() {
	fmt.Println("--- ex12Enums ---")
	fmt.Println(RoleAdmin)
	fmt.Println(RoleUser)
	fmt.Println(RoleGuest)
	fmt.Println(json.Marshal(RoleAdmin))

	fmt.Println(StatusActive)
	fmt.Println(StatusInactive)
	fmt.Println(json.Marshal(StatusActive))

	p := PermRead | PermWrite
	fmt.Println(p)                       // 3
	fmt.Println(p & PermWrite != 0)      // true — bit is set
	fmt.Println(p & PermDelete != 0)     // false

}
