# BettErr

BettErr is a Go library that enhances error handling by providing stack traces and multiple error formatting styles similar to other programming languages.

## Why BettErr?
- Better Debug Information: Get stack traces for easier debugging
- Flexible Formatting: Choose the error format that best suits your needs (Go-like, Java-like, and JSON)
- Simple API: Intuitive methods for creating and handling errors
- Go Compatibility: Works seamlessly with standard Go errors

## Roadmap
- Put runnable examples in the doc
- Add a way to disable stacktrace (setting + env ?)

## Installation

```bash
go get github.com/jjunac/betterr
```

## Usage

```go
// Create a new error with stack trace
err := betterr.New("something went wrong")

// Wrap an existing error with stack trace
plainErr := returnsAnError()
wrappedErr := betterr.Wrap(plainErr)

// Decorate an error with additional context
decoratedErr := betterr.Decorate(err, "failed to process")
// Or with formatting
decoratedErr = betterr.Decoratef(err, "failed to process item %d", 123)
```

## Formatting Errors

BettErr supports multiple formatting styles. The `Error()` methods of the error use the default formatter (Java style by default). \
You can set the default formatter by changing `betterr.DefaultFortmatter`:
```go
betterr.DefaultFortmatter = &betterr.JavaStyleFormatter{}
```

### Go Style

```go
fmt.Println(formatter.Format(&betterr.GoStyleFormatter{}))
// Output: failed to process: something went wrong
```

### Java Style

```go
fmt.Println(formatter.Format(&betterr.JavaStyleFormatter{}))
// Output:
// failed to process
//     at github.com/myapp.MyFunction (file.go:123)
//     at github.com/myapp.main (main.go:45)
// Caused by: something went wrong
//     at github.com/myapp.OtherFunction (file.go:100)
```

### JSON (useful for monitoring for instance)

```go
fmt.Println(formatter.Format(&betterr.JsonFormatter{}))
// Output:
// {
//     "message": "failed to process",
//     "stack": [
//         {
//             "file": "file.go",
//             "function": "github.com/myapp.MyFunction",
//             "line": 123
//         },
//         {
//             "file": "main.go",
//             "function": "github.com/myapp.main",
//             "line": 45
//         }
//     ],
//     "cause": {
//         "message": "something went wrong",
//         "stack": [
//             {
//                 "file": "file.go",
//                 "function": "github.com/myapp.OtherFunction",
//                 "line": 42
//             }
//         ]
//     }
// }
```

## Benchmark

BettErr is faster than the other error handling libraries (such as [eris](https://github.com/rotisserie/eris) and [errorx](https://github.com/joomcode/errorx)),
but obviously slower than standard's Go errors, since they don't offer the stack trace feature (and runtime stack unwinding takes time).

goos: darwin
goarch: arm64
pkg: github.com/jjunac/betterr/benchmark
cpu: Apple M3 Pro
|      Benchmark      |    Runs     | Time per error creation |
|---------------------|-------------|-------------------------|
| 10-frame stack      |             |                         |
| ---Errors           |    53471764 |             22.65 ns/op |
| **---Betterr**      | **1862188** |          **614.2 ns/op**| 
| ---Errorx           |     1578981 |             853.7 ns/op |
| ---Eris             |      444178 |              2824 ns/op |
| 100-frame stack     |             |                         |
| ---Errors           |     2480102 |             479.0 ns/op |
| **---Betterr**      |  **938510** |          **1302 ns/op** |
| ---Errorx           |      401744 |              3017 ns/op |
| ---Eris             |      134414 |              8903 ns/op |

