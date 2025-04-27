# DAPA Backend
This repository contains the backend of the DAPA application developed using **Go** with the **Gin** framework, as well as its setup and execution using **Docker**.

## Getting started
The project includes a `docker-compose.yml` file to run the application and its database with Docker.

### Prerequisites
- Git
- Docker
- Docker Compose

### Installation
1. **Clone** both the backend and database repositories. 
```bash
git clone https://github.com/JuanDsm04/dapa-backend.git
git clone https://github.com/vicperezch/dapa-database.git
cd dapa-backend
```

2. **Create** the `.env` file in the root of the project.
```env
POSTGRES_USER=youruser
POSTGRES_DB=yourdb
POSTGRES_PASSWORD=yourpassword
JWT_SECRET=yoursecret
```

3. **Build** and serve using Docker. By default, the server will run on `http://localhost:8080`.
```bash
docker-compose up --build
```

4. **Stop** the containers.
```bash
docker-compose down
```

## Usage
The API uses **Swagger** to provide documentation. To access it and view the available endpoints go to `http://localhost:8080/swagger/index.html` after running the containers.
