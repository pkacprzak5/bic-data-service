FROM golang:1.23.4-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/api ./cmd

EXPOSE 8080

CMD ["./bin/api"]