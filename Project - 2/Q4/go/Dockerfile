FROM golang:1.24-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o q4server server.go
EXPOSE 6001 6002 6003 6004 6005
CMD ["./q4server","1","6001", "q4node2:6002,q4node3:6003,q4node4:6004,q4node5:6005"]