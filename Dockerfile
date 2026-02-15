FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE
RUN go build -o service ./cmd/${SERVICE}

EXPOSE 50051 50052 50053 8080

CMD ["./service"]
