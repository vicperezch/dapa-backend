# DAPA Backend

This repository contains the backend of the DAPA application developed using **Go** with the **Gin** framework, as well as its setup and execution using **Docker**.

## Getting started

1. Clone the repository 

```bash
git clone https://github.com/JuanDsm04/dapa-backend
cd dapa-backend
```

2. Install dependencies

```bash
go mod tidy
```

## Project Setup

1. Make sure to configure your `.env` file if needed. Example:

```bash
PORT=8080
```

By default, the server will run on `http://localhost:8080`.

2. Run the project locally

```bash
go run src/cmd/main.go
```

To test if the backend is working correctly, you can add a basic test endpoint in `main.go`:

```go
router.GET("/api/ping", func(c *gin.Context) {
  c.JSON(200, gin.H{"message": "pong"})
})
```

Then visit:

```
http://localhost:8080/api/ping
```

You should receive the following response:

```json
{"message": "pong"}
```

### Build and serve through Docker

```bash
docker compose up --build
```

### Stop Docker Compose

```bash
docker compose down
```
