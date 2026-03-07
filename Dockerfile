# ---------- Build Stage ----------
FROM golang:1.25 AS builder

WORKDIR /app

# download dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy source code
COPY . .

# build server binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o server ./cmd/server

# ---------- Runtime Stage ----------
FROM alpine:3.19

WORKDIR /app

# add certificates
RUN apk add --no-cache ca-certificates

# copy binary from builder
COPY --from=builder /app/server .

# copy configuration and data
COPY config.yaml .
COPY data ./data

# server port
EXPOSE 8080

# start server
ENTRYPOINT ["./server"]