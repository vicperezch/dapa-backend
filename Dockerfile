FROM golang:1.24

WORKDIR /app

# Instalar swaggo
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY /go.* ./
RUN go mod download

COPY ./ .
COPY .env ./

# Generar documentación
RUN swag init -g ./cmd/main.go -o ./docs
RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]
