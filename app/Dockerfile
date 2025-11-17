FROM golang:1.25-alpine AS builder

WORKDIR /app

# only copy go.mod first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# copy the rest of the app
COPY . .

# build the binary - static linking for scratch
RUN CGO_ENABLED=0 go build \
    -trimpath \
    -ldflags "-s -w -X main.version=${VERSION}" \
    -o /tunnellm

# --- Runtime stage ---
FROM scratch

COPY --from=builder /tunnellm /tunnellm
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/tunnellm"]
