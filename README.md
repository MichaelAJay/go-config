# Go Config

A lightweight, flexible configuration management package for Go applications.

## Features

- Multiple configuration sources (YAML, JSON, environment variables)
- Type-safe configuration access
- Configuration validation
- Default values
- Hierarchical configuration with overrides
- Thread-safe operations

## Installation

```bash
go get github.com/MichaelAJay/go-config
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/MichaelAJay/go-config"
)

func main() {
    // Create a new configuration manager
    cfg := config.New()

    // Load from a YAML file
    if err := cfg.Load(&config.FileSource{Path: "config.yaml"}); err != nil {
        log.Fatal(err)
    }

    // Load from environment variables
    if err := cfg.Load(&config.EnvSource{Prefix: "APP_"}); err != nil {
        log.Fatal(err)
    }

    // Set default values
    defaults := &config.DefaultSource{
        Values: map[string]interface{}{
            "port":     8080,
            "host":     "localhost",
            "debug":    false,
            "timeout":  30,
        },
    }
    if err := cfg.Load(defaults); err != nil {
        log.Fatal(err)
    }

    // Add validators
    cfg.AddValidator(&config.RequiredValidator{
        Keys: []string{"port", "host"},
    })

    cfg.AddValidator(&config.RangeValidator{
        Key:   "port",
        Min:   1024,
        Max:   65535,
        IsInt: true,
    })

    // Validate configuration
    if err := cfg.Validate(); err != nil {
        log.Fatal(err)
    }

    // Access configuration values
    port, _ := cfg.GetInt("port")
    host, _ := cfg.GetString("host")
    debug, _ := cfg.GetBool("debug")
    timeout, _ := cfg.GetInt("timeout")

    fmt.Printf("Server running on %s:%d (debug: %v, timeout: %d)\n", host, port, debug, timeout)
}
```

### Configuration File Example (config.yaml)

```yaml
server:
  host: localhost
  port: 8080
  timeout: 30
  debug: false

database:
  host: localhost
  port: 5432
  name: myapp
  user: postgres
  password: secret

logging:
  level: info
  format: json
  output: stdout
```

### Environment Variables

Environment variables can override configuration values. They should be prefixed with `APP_` (or your chosen prefix) and use underscores for hierarchy:

```bash
APP_SERVER_PORT=9090
APP_DATABASE_PASSWORD=newsecret
APP_LOGGING_LEVEL=debug
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.