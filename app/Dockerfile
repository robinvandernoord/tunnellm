FROM golang:1.25-alpine AS builder

WORKDIR /app

# only copy go.mod first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# copy the rest of the app
COPY . .

# build the binary
RUN go build \
    -trimpath \
    -ldflags "-s -w -X main.version=${VERSION}" \
    -o /tunnellm

# --- Runtime stage ---
FROM alpine:3.21

# create a non-root user (optional but good practice)
RUN adduser -D -g '' appuser

# copy the built binary from builder
COPY --from=builder /tunnellm /usr/local/bin/tunnellm

USER appuser

ENTRYPOINT ["/usr/local/bin/tunnellm"]
