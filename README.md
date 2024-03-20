# Auth Proxy

**English** | [简体中文](./README.zh-CN.md)

## Overview

A minimalist reverse proxy service for web service authentication

- Minimalist UI
- Containerized deployment

## Usage

You can download the appropriate software package for your platform from [releases](https://github.com/wengchaoxi/auth-proxy/releases/latest) and unzip it, edit the `.env.example` file, rename it to `.env`, and then execute the binary file.

```yml
# Runtime Configuration
HOST = "0.0.0.0"
PORT = 18000

# Real service address (User -> Auth Proxy -> Real service address)
TARGET_URL = "http://localhost:8000"

# Authentication access key, default is `whoami`
AUTH_ACCESS_KEY = "whoami"

# Authentication validity period, need to re-enter the access key after expiration, default is `24h`
AUTH_EXPIRATION = "24h"
```

Or, you can use Docker:

```
docker run --rm -p 18000:18000 -e TARGET_URL=http://localhost:8000 -e AUTH_ACCESS_KEY=whoami wengchaoxi/auth-proxy:latest
```

Then visit: http://localhost:18000

## Example

> Add access authentication for the `traefik/whoami` service, with the access key as `whoami`

Edit docker-compose.yml and run `docker compose up -d`, then visit: http://localhost:18000

```yml
version: '3'

services:
  proxy:
    image: wengchaoxi/auth-proxy:latest
    ports:
      - 18000:18000
    environment:
      - HOST=0.0.0.0
      - PORT=18000
      - TARGET_URL=http://whoami:8000
      - AUTH_ACCESS_KEY=whoami
      - AUTH_EXPIRATION=24h

  whoami:
    image: traefik/whoami
    command:
      - --port=8000
    ports:
      - "8000:8000"
```
