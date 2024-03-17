# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS builder
RUN apk --no-cache --no-progress add ca-certificates tzdata \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

WORKDIR /src
COPY . .
RUN --mount=type=cache,target=/go --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -a --trimpath --ldflags='-s -w' -o auth-proxy .

FROM scratch
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app
COPY --from=builder /src/auth-proxy ./
COPY --from=builder /src/web ./web

EXPOSE 18000
ENTRYPOINT ["/app/auth-proxy"]
