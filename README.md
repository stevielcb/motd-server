# motd-server

- [motd-server](#motd-server)
  - [Overview](#overview)
  - [Features](#features)
  - [How It Works](#how-it-works)
  - [Configuration](#configuration)
  - [Running](#running)
  - [Notes](#notes)
  - [License](#license)

## Overview

A lightweight TCP server that serves a random file from a specified cache directory when a client connects.

## Features

- Simple TCP server.
- Randomly selects and sends a cached message (file) to each connecting client.
- Configurable through environment variables.

## How It Works

- Listens on the configured host and port (`MOTD_LISTEN_HOST` and `MOTD_LISTEN_PORT`).
- Upon each client connection, selects a random file from the cache directory (`MOTD_CACHE_DIR`).
- Sends the contents of the file to the client, then closes the connection.

## Configuration

Environment variables:

| Variable                  | Default         | Description                                    |
|----------------------------|-----------------|------------------------------------------------|
| MOTD_LISTEN_HOST           | localhost       | Host address to bind the server.               |
| MOTD_LISTEN_PORT           | 4200            | Port to listen on.                             |
| MOTD_CACHE_DIR             | ~/.motd         | Directory containing cached message files.    |
| MOTD_GIPHY_API_KEY_FILE    | ~/.giphy-api    | File containing Giphy API Key (optional).      |
| MOTD_DOWNLOAD_INTERVAL     | 10              | Interval for downloading new files (optional). |
| MOTD_CLEANUP_INTERVAL      | 60              | Interval for cache cleanup (optional).         |
| MOTD_GIPHY_TAGS            | (none)          | Giphy tags for selecting GIFs (optional).      |

## Running

1. Build the server:

   ```bash
   go build -o motd-server
   ```

2. Run the server:

   ```bash
   ./motd-server
   ```

## Notes

- The server will automatically create the cache directory if it does not exist.
- It expects message files already to be present in the cache directory unless extended with download functionality.

## License

MIT License
