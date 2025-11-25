FROM golang:1.22-alpine AS builder

WORKDIR /

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target="/root/.cache/go-build" \
    go build -o /bin/pr -trimpath -ldflags "-s -w" ./cmd/

FROM alpine:latest

WORKDIR /app

COPY --from=builder /bin/pr .

CMD ["./pr"]
