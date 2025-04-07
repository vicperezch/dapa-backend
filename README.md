# DAPA Backend

Este repositorio contiene la aplicación backend del proyecto **DAPA**, desarrollada en **Go (Golang)** utilizando el framework **Gin**, y preparada para ejecutarse mediante **Docker**.

## Comenzando 

1. Clona el repositorio:

```bash
git clone https://github.com/JuanDsm04/dapa-backend
cd dapa-backend
```

2. Instala las dependencias necesarias:

```bash
go mod tidy
```

## Configuración del Proyecto 

1. Asegúrate de tener configurado el archivo `.env` si es necesario. Un ejemplo puede ser:

```bash
PORT=8080
```

Por defecto, el servidor se ejecuta en `http://localhost:8080`.

2. Para ejecutar el servidor localmente:

```bash
go run src/cmd/main.go
```

3. Puedes agregar un endpoint de prueba para verificar que todo funcione correctamente. Añade lo siguiente en tu archivo `main.go`:

```go
router.GET("/api/ping", func(c *gin.Context) {
  c.JSON(200, gin.H{"message": "pong"})
})
```

Luego accede a:

```
http://localhost:8080/api/ping
```

Y deberías ver una respuesta como esta:

```json
{"message": "pong"}
```

## Docker 

### Construir y levantar el contenedor con Docker

```bash
docker compose up --build
```

### Detener los servicios de Docker Compose

```bash
docker compose down
```

## Endpoints disponibles

Actualmente disponibles:

- `POST /api/user`: Crear un nuevo usuario.
- `GET /api/user`: Obtener lista de usuarios.
- `GET /api/ping`: (Si fue añadido) Endpoint de prueba para verificar conexión.

## Notas adicionales

- Este backend está diseñado para ser consumido por el frontend del proyecto DAPA.
- Asegúrate de que tanto el frontend como el backend estén utilizando la misma URL base para facilitar la comunicación entre ambos servicios.
- Si usas Docker, asegúrate de que los servicios estén correctamente configurados en `docker-compose.yml`.

