# -------------------
# Stage 1: Build Go (for DC)
# -------------------
    FROM golang:1.24 AS go-builder
    WORKDIR /app
    
    COPY go.mod go.sum ./
    RUN go mod download
    
    COPY . .
    
    RUN CGO_ENABLED=0 GOOS=linux go build -o /app/dissemination dissemination.go
    
    # -------------------
    # Stage 2: Final image (Python FD and Go DC binary)
    # -------------------
    FROM python:3.9
    WORKDIR /app
    
    COPY --from=go-builder /app/dissemination /usr/local/bin/dissemination
    
    COPY . .
    
    RUN python -m pip install --upgrade pip
    RUN pip install --no-cache-dir grpcio grpcio-tools
    
    CMD ["sh", "-c", "dissemination & python main.py"]
    