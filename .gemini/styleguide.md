# Golang Code Review Styleguide

Please review the following Go code strictly against standard Go idioms, the Uber Go Style Guide, and the Google Go Style Guide. Flag any code that violates these principles and provide the idiomatic correction.

## 1. Code Readability & Philosophy

* **Optimize for the Reader:** Code must prioritize readability over concise "cleverness". Avoid overly complex one-liners if multiple lines are easier to read.
* **Line of Sight:** Keep the "happy path" aligned to the left edge. Avoid deeply nested `if` statements. Use guard clauses and early returns to handle errors and edge cases immediately. Never use `else` if the `if` block ends in a `return`, `break`, or `continue`.
* **Comments Explain "Why":** Code should explain *what* it is doing; comments should explain *why*. Flag complex logic that lacks a comment explaining the underlying business reason or edge case.

## 2. Pointers and Interfaces

* **No Pointers to Interfaces:** Never use a pointer to an interface (e.g., `*error`, `*io.Reader`). Interfaces are already reference types.
* **Receiver Types:** Use value receivers by default. Use pointer receivers only if the method needs to mutate the receiver, if the struct contains a `sync.Mutex`, or if the struct is very large.

## 3. Naming Conventions

* **Scope-Based Length:** Short variables (e.g., `i`, `c`, `err`) are acceptable only for small, highly localized scopes. For larger scopes or package-level variables, names must be fully descriptive.
* **Receiver Names:** Receiver names must be short (1 or 2 letters) and consistent across all methods of a type. Never use `this`, `self`, or `me`.
* **Errors:** Error variables must be prefixed with `err` (e.g., `errNotFound`). Custom error types must end with `Error` (e.g., `NotFoundError`).
* **Unexported Globals:** Unexported package-level variables should be prefixed with an underscore (e.g., `_defaultTimeout`).

## 4. Error Handling and Panics

* **Error Wrapping:** Always wrap errors with context using `fmt.Errorf("...: %w", err)` when passing them up the stack.
* **Error Messages:** Error strings must not be capitalized (unless beginning with a proper noun or acronym) and must not end with punctuation (e.g., `fmt.Errorf("failed to load file: %w", err)`).
* **Type Assertions:** Always use the "comma ok" idiom for type assertions to prevent panics (e.g., `t, ok := i.(string)`).
* **No Panics:** Do not use `panic` in normal code. Panics are only acceptable in `init()` functions or during application startup if a critical dependency is missing.

## 5. Concurrency and Synchronization

* **Zero-value Mutex:** Use the zero-value of `sync.Mutex` and `sync.RWMutex`. Do not use pointers to mutexes unless they are embedded in a struct that is also passed by pointer.
* **Channels:** Channel sizes should typically be 1 or unbuffered. If a larger buffer is used, ensure there is a comment explaining why the specific size was chosen.
* **Goroutine Lifetimes:** Ensure that any started goroutine has a clear exit path. Flag potential goroutine leaks.

## 6. Performance and Initialization

* **String Formatting:** Prefer `strconv` (e.g., `strconv.Itoa`) over `fmt` (e.g., `fmt.Sprintf`) for converting base types to strings.
* **Preallocate Collections:** Always specify capacity hints when making maps or slices if the size is known: `make(map[T1]T2, hint)` or `make([]T, 0, capacity)`.
* **Avoid `init()`:** Avoid using `init()` functions to set up state or dependencies. Prefer explicit `Setup()` or `New...()` constructor functions.

## 7. Testing

* **Table-Driven Tests:** Enforce the use of table-driven tests (`[]struct{ name string ... }`) for repetitive test cases.
