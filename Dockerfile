FROM golang:1.24

WORKDIR /app

COPY /go.* ./
RUN go mod download

COPY ./ .
COPY .env ./

RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]
