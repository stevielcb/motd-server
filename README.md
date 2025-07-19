# motd-server

- [motd-server](#motd-server)
  - [Overview](#overview)
  - [Features](#features)
  - [Architecture](#architecture)
    - [Key Design Principles](#key-design-principles)
  - [How It Works](#how-it-works)
  - [Configuration](#configuration)
  - [Running](#running)
  - [Development](#development)
    - [Building](#building)
    - [Testing](#testing)
    - [Adding New Services](#adding-new-services)
    - [Project Structure](#project-structure)
  - [License](#license)

## Overview

A lightweight TCP server that serves a random file from a specified cache directory when a client connects. The server automatically downloads content from external services (Giphy and XKCD) and caches it for serving.

## Features

- Simple TCP server with graceful shutdown
- Automatic content downloading from Giphy and XKCD APIs
- Intelligent cache management with size limits
- Configurable through environment variables
- Clean, testable architecture with dependency injection
- Comprehensive error handling and logging

## Architecture

The application follows a clean architecture pattern with clear separation of responsibilities:

```
motd-server/
├── app/                    # Application container and lifecycle
├── internal/
│   ├── cache/             # Cache management operations
│   ├── config/            # Configuration loading and validation
│   ├── server/            # TCP server implementation
│   └── services/          # External service integrations
│       ├── giphy/         # Giphy API client
│       └── xkcd/          # XKCD API client
├── main.go                # Entry point with graceful shutdown
└── README.md              # This file
```

### Key Design Principles

- **Dependency Injection**: All dependencies are explicitly passed through constructors
- **Interface-Based Design**: Services use interfaces for loose coupling and testability
- **Single Responsibility**: Each package has a clear, focused purpose
- **Error Handling**: Proper error propagation with context
- **Graceful Shutdown**: Signal handling with proper cleanup

## How It Works

1. **Startup**: The application loads configuration, initializes all services, and starts background workers
2. **Content Download**: Background workers periodically fetch new content from Giphy and XKCD APIs
3. **Caching**: Downloaded content is stored in the local cache directory with metadata
4. **Serving**: When clients connect, the server randomly selects and serves cached content
5. **Cleanup**: Background workers periodically clean up old cache files to maintain size limits

## Configuration

Environment variables:

| Variable                  | Default         | Description                                    |
|----------------------------|-----------------|------------------------------------------------|
| MOTD_LISTEN_HOST           | localhost       | Host address to bind the server.               |
| MOTD_LISTEN_PORT           | 4200            | Port to listen on.                             |
| MOTD_CACHE_DIR             | ~/.motd         | Directory containing cached message files.    |
| MOTD_GIPHY_API_KEY_FILE    | ~/.giphy-api    | File containing Giphy API Key (optional).      |
| MOTD_DOWNLOAD_INTERVAL     | 10              | Interval for downloading new files (seconds).  |
| MOTD_CLEANUP_INTERVAL      | 60              | Interval for cache cleanup (seconds).          |
| MOTD_GIPHY_TAGS            | (none)          | Giphy tags for selecting GIFs (optional).      |
| MOTD_CACHE_MAX_FILES       | 50              | Maximum number of cached files to keep.        |

## Running

1. Build the server:

   ```bash
   go build -o motd-server
   ```

2. Run the server:

   ```bash
   ./motd-server
   ```

3. Connect to the server:

   ```bash
   telnet localhost 4200
   ```

## Development

### Building

```bash
go build -o motd-server
```

### Testing

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/services/...
```

## CI/CD

This project uses GitHub Actions for continuous integration and deployment:

### Pull Request Testing

When a pull request is submitted, the following checks are automatically run:

- **Go Tests**: Runs all unit tests with race detection and coverage reporting
- **Code Quality**: Runs golangci-lint for code quality checks
- **Build Verification**: Ensures the application builds successfully
- **Multi-platform Build**: Tests building for Linux, macOS, and Windows

### Automated Releases

When code is merged to the `main` branch:

1. **Auto-tagging**: Automatically creates a new patch version tag (e.g., v1.0.1)
2. **Release Creation**: Creates a GitHub release with the new tag
3. **Binary Distribution**: Builds and uploads binaries for multiple platforms:
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64)
4. **Checksums**: Provides SHA256 checksums for all binaries

### Version Information

The application includes version information that can be displayed with:

```bash
./motd-server -version
```

Version is automatically set during the build process from Git tags.

### Adding New Services

To add a new MOTD provider:

1. Create a new service in `internal/services/`
2. Implement the appropriate interface
3. Add the service to the `services.Manager`
4. Update configuration as needed

Example:

```go
type MyService struct {
    // service implementation
}

func (s *MyService) GetRandom() (string, error) {
    // fetch content from your service
    return "https://example.com/content", nil
}
```

### Project Structure

- **`app/`**: Application lifecycle and dependency management
- **`internal/config/`**: Configuration loading and validation
- **`internal/cache/`**: Cache operations and file management
- **`internal/server/`**: TCP server implementation
- **`internal/services/`**: External service integrations
  - **`giphy/`**: Giphy API client
  - **`xkcd/`**: XKCD API client

## License

MIT License
