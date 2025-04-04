FROM golang:1.24

WORKDIR /app

COPY /src/go.* ./
RUN go mod download

COPY . .

RUN go build -o main ./src/cmd/main.go

EXPOSE 8080

CMD ["./main"]
