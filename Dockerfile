FROM golang:1.24-alpine AS builder
WORKDIR /app

# Instala dependencias necesarias para build
RUN apk add --no-cache git

# Copiar mod y descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Generar la documentación Swagger
RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init -g ./cmd/main.go -o ./docs

# Compilar el binario para producción
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

FROM alpine:3.20

# Crear un usuario no root
RUN adduser -D -u 1000 appuser
USER appuser

WORKDIR /app

# Copiar el binario y archivos necesarios
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs
COPY .env .env

# Exponer puerto del backend
EXPOSE 8080

CMD ["./main"]
