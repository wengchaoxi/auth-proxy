# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS builder
RUN apk --no-cache --no-progress add ca-certificates tzdata \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

FROM scratch
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app
COPY auth-proxy ./
COPY web ./web
COPY LICENSE ./

EXPOSE 18000
ENTRYPOINT ["/app/auth-proxy"]
