# -------------------
#   Stage 1: Build Go
# -------------------
    FROM golang:1.24 AS go-builder
    WORKDIR /app
    
    COPY go.mod go.sum ./
    RUN go mod download
    
    COPY . .
    RUN CGO_ENABLED=0 GOOS=linux go build -o /app/dissemination dissemination.go
    
    # ----------------------------
    #   Stage 2: Final image
    # ----------------------------
    FROM python:3.9
    
    WORKDIR /app
    
    # Copy the compiled Go binary
    COPY --from=go-builder /app/dissemination /usr/local/bin/dissemination
    
    # Copy python code, proto, etc.
    COPY . /app
    
    # Upgrade pip (optional but recommended) then install grpcio
    RUN python -m pip install --upgrade pip
    RUN pip install --no-cache-dir grpcio grpcio-tools
    
    # If you have other dependencies, install them here:
    # RUN pip install --no-cache-dir -r requirements.txt
    
    EXPOSE 50051 50052 50053 50054 50055 6000
    
    CMD ["sh", "-c", "dissemination & python main.py"]
    