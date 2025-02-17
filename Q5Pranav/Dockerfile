FROM golang:1.16

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o dissemination-component

EXPOSE 50051

CMD ["./dissemination-component"]
