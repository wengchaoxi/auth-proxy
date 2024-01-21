FROM golang:1.22-alpine AS builder

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct

WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/go --mount=type=cache,target=/root/.cache/go-build \
    go build -o ./auth-proxy .

FROM alpine

WORKDIR /app
COPY --from=builder /app/auth-proxy ./
COPY --from=builder /app/web ./

EXPOSE 8080
CMD ["/app/auth-proxy"]
