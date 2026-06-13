# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w' -o server ./cmd/server

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

RUN addgroup -S still && adduser -S still -G still
USER still

WORKDIR /app

COPY --from=builder /build/server ./server

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./server"]
