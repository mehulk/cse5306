# Dockerfile for the Go-based Q4 LogReplicator
FROM golang:1.24-alpine


WORKDIR /app

# Copy go.mod / go.sum first if you have them, then do go mod download
COPY go.mod go.sum ./
RUN go mod download

# Copy your entire code
COPY . .

# Build the Go server
RUN go build -o q4server server.go

# Expose relevant port for gRPC
EXPOSE 6001

# Start the Go server
CMD ["./q4server", "1", "6001", "peer2:6002,peer3:6003,peer4:6004,peer5:6005"]
