# Goticks

[![CI](https://github.com/parametalol/goticks/actions/workflows/ci.yml/badge.svg)](https://github.com/parametalol/goticks/actions/workflows/ci.yml) [![codecov](https://codecov.io/gh/parametalol/goticks/branch/main/graph/badge.svg)](https://codecov.io/gh/parametalol/goticks) [![Go Reference](https://pkg.go.dev/badge/github.com/parametalol/goticks.svg)](https://pkg.go.dev/github.com/parametalol/goticks) [![Go Report Card](https://goreportcard.com/badge/github.com/parametalol/goticks)](https://goreportcard.com/report/github.com/parametalol/goticks) [![License](https://img.shields.io/github/license/parametalol/goticks)](./LICENSE)

Goticks is a lightweight Go library for building and managing periodic tasks with support for cancellable contexts, customizable tickers, error handling, and more.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [API Reference](#api-reference)
- [Contributing](#contributing)
- [License](#license)

## Features

- Modular design: clear separation between tasks, tickers, and runners.
- Cancellable contexts for graceful shutdown.
- Customizable tick generators that are restartable and stoppable.
- Immediate first tick on start.
- Built-in retry and error handling support.
- Zero dependencies outside the Go standard library.

## Installation

```bash
go get github.com/parametalol/goticks
```

## Usage

```go
package main

import (
    "fmt"
    "time"

    "github.com/parametalol/goticks"
    "github.com/parametalol/goticks/ticker"
)

func main() {
    start := time.Now()
    t := ticker.NewTimer(time.Second)
    goticks.NewTask(t, func(t time.Time) {
        fmt.Println("Current time:", t.Sub(start).Round(time.Second))
    }).Start()

    // Let it run for 3 seconds.
    time.Sleep(3 * time.Second)
    t.Stop()
}
```

## API Reference

Detailed documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/parametalol/goticks).

## Contributing

Contributions, issues, and feature requests are welcome! Please check [issues](https://github.com/parametalol/goticks/issues) or submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
