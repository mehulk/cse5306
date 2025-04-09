# Dockerfile for "q4server"
FROM golang:1.24-alpine AS build

WORKDIR /app

# Copy mod files first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest
COPY . .

# Build the binary
RUN go build -o /q4server server.go

# Final stage
FROM alpine:3.18
COPY --from=build /q4server /q4server
EXPOSE 6001
ENTRYPOINT ["/q4server"]
