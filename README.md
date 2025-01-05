# `slogger`

`slogger` is a lightweight wrapper for Go's standard logging library [slog](https://pkg.go.dev/log/slog). It simplifies logging with an intuitive API, enhanced context support, and flexible configuration options.

---

## Features

- **Simplified API**  
  Log messages with concise and intuitive function calls.
- **Context Support**  
  Add custom fields to logs using `context.Context`.
- **Source Information**  
  Automatically appends the file name and line number to each log message.
- **Flexible Configuration**  
  Customize log levels, output destinations, path trimming, and additional context fields.
- **Standard Logger Compatibility**  
  Fully compatible with `log.Logger` for seamless integration with existing libraries.

---

## Installation

Install the library using:

```bash
go get github.com/mickamy/slogger
```

---

## Usage

### Initialization

Before using `slogger`, initialize it with a configuration.

```go
package main

import (
    "context"
	"github.com/mickamy/slogger"
)

func main() {
	// Initialize slogger
	slogger.Init(slogger.Config{
		Level:          slogger.LevelInfo,
		TrimPathPrefix: "/path/to/the/project/",
		ContextFieldsExtractor: func(ctx context.Context) []any {
			return []any{"userID", ctx.Value("userID")}
		},
	})

	// Log a message
	slogger.Info("Application started")
}
```

---

### Logging Levels

`slogger` supports multiple logging levels.

```go
slogger.Debug("This is a debug message")
slogger.Info("This is an info message")
slogger.Warn("This is a warning message")
slogger.Error("This is an error message")
```

---

### Context-Aware Logging

Add custom fields to your logs using `context.Context`.

```go
ctx := context.WithValue(context.Background(), "userID", 12345)
slogger.InfoCtx(ctx, "User logged in")
```

Example output:

```json
{
  "level": "INFO",
  "msg": "User logged in",
  "source": "main.go:42",
  "userID": 12345
}
```

---

### Using as a Standard Logger

You can use `slogger` as a `log.Logger` for compatibility with other existing libraries.

```go 
stdLogger := slogger.StandardLogger(slogger.LevelInfo)
stdLogger.Println("This is a standard log message")
```
---

## Configuration Options

You can customize `slogger` by providing a Config struct during initialization.

| Field                    | Description                                                              | Default                  |
|--------------------------|--------------------------------------------------------------------------|--------------------------|
| `Level`                  | Minimum log level (`LevelDebug`, `LevelInfo`, `LevelWarn`, `LevelError`) | `LevelInfo`              |
| `Outputs`                | Log output destinations (`io.Writer`)                                    | `[]io.Writer{os.Stdout}` |
| `TrimPathPrefix`         | Prefix to trim from file paths in log output                             | Empty string             |
| `ContextFieldsExtractor` | Function to extract additional fields from `context.Context`             | `nil` (none)             |

---

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
