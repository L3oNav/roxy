# Roxy - Lightweight Reverse Proxy Server

**Roxy** is a lightweight and highly configurable reverse proxy server written in Go. It supports load balancing with algorithms like Weighted Round Robin (WRR) and offers an easy-to-use configuration system.

## Features

- **Reverse Proxy:** Route incoming HTTP requests to backend servers.
- **Load Balancing:** Distribute traffic across multiple backend servers using the Weighted Round Robin (WRR) algorithm.
- **Configurable Logging:** Customizable logging levels and output formats.
- **Graceful Shutdown:** Ensure existing connections are properly handled during shutdown.
- **File Serving:** Serve static files directly from the server.
- **Configurable through TOML:** Easily define server behavior and routing rules using a TOML configuration file.

## Getting Started

### Prerequisites

- Go 1.20 or later

### Installation

1. Clone the repository:

    ```bash
    git clone git@github.com:L3oNav/roxy.git
    cd roxy
    ```

2. Build the project:

    ```bash
    go build -o roxy ./src/main.go
    ```

3. Create a configuration file `config.toml` based on the example below.

### Configuration

Roxy uses a TOML configuration file to define server settings and routing rules. Below is an example configuration:

```toml
[server]
name = "roxy"
logfile = "logs/access.log"
loglevel = "debug"
max_connections = 1024
listen = ["127.0.0.1:8100", "192.168.1.2:8100"]

[[match]]
uri = "/"
serve = "/static"

[[match]]
uri = "/api"
algorithm = "WRR"
forward = [
    { address = "127.0.0.1:8080", weight = 1 },
    { address = "127.0.0.1:8081", weight = 3 },
    { address = "127.0.0.1:8082", weight = 2 },
]

```

### Usage

Run the proxy server with:

```bash
./roxy
```

By default, the server listens on the address specified in the `config.toml` file.

### Running Tests

To run the tests:

```bash
go test ./src/...
```

This will execute all the test files, including the configuration loader and other components.

### Contributing

Contributions are welcome! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new feature branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -am 'Add some feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Open a pull request.

Please ensure your code adheres to the project's coding standards and passes all tests.

### License

This project is licensed under the MIT License
