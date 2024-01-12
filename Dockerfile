FROM golang:1.21-alpine AS builder

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/go --mount=type=cache,target=/root/.cache/go-build \
    go build -o ./proxy .

FROM alpine

WORKDIR /app
COPY --from=builder /app/proxy ./
COPY --from=builder /app/web .

EXPOSE 8080
CMD ["/app/proxy"]
