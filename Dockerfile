FROM golang:1.24-alpine AS builder
WORKDIR /app

# Instala dependencias necesarias para build
RUN apk add --no-cache git

# Copiar mod y descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar el binario para producción
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go


FROM alpine:3.20

# Crear un usuario no root
RUN adduser -D -u 1000 appuser

WORKDIR /app

# Copiar el binario y archivos necesarios
COPY --from=builder /app/main .
COPY .env .env

USER appuser

# Exponer puerto del backend
EXPOSE 8080

CMD ["./main"]
