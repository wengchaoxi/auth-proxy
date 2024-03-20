# Auth Proxy

**简体中文** | [English](./README.md)

## 概述

一个极简的用于 Web 服务鉴权的反向代理服务

- 极简的 UI
- 容器化部署

## 使用说明

你可以在 [发布页面](https://github.com/wengchaoxi/auth-proxy/releases/latest) 下载相应平台的软件包并解压，编辑其中 `.env.example` 文件，并重命名为 `.env`，然后执行二进制文件
```yml
# 运行配置
HOST = "0.0.0.0"
PORT = 18000

# 真实服务地址（用户 -> Auth Proxy -> 真实服务地址）
TARGET_URL = "http://localhost:8000"

# 认证访问密钥，默认 `whoami`
AUTH_ACCESS_KEY = "whoami"

# 认证有效期，过期后需要重新输入访问密钥，默认 `24h`
AUTH_EXPIRATION = "24h"
```

或者，你可以使用 Docker

```sh
docker run --rm -p 18000:18000 -e TARGET_URL=http://localhost:8000 -e AUTH_ACCESS_KEY=whoami wengchaoxi/auth-proxy:latest
```

然后访问：http://localhost:18000

## 部署示例

> 为 `traefik/whoami` 服务添加访问鉴权，访问密钥为 `whoami`


编辑 docker-compose.yml 并运行 `docker compose up -d`，然后访问：http://localhost:18000

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
