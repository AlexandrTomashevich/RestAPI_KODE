FROM golang:1.21.1

WORKDIR /app

COPY go.mod go.sum ./

RUN go build -o main .

CMD ["./main"]