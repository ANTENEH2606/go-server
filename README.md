# Album REST API Server with Go and PostgreSQL

This is a simple RESTful API server written in Go that manages album records stored in a PostgreSQL database. It supports basic CRUD operations (Create, Read, Delete) on album entities via HTTP endpoints.

## Features

- Connects to PostgreSQL using go-pg ORM
- CRUD endpoints for `albums` resource:
  - `GET /albums` — list all albums
  - `POST /albums` — create a new album
  - `GET /albums/{id}` — get album by ID
  - `DELETE /albums/{id}` — delete album by ID
- Uses environment variables for configuration
- JSON request and response format
- Basic error handling with JSON error responses

## Album Model

The album resource contains the following fields:

| Field  | Type    | Description         |
|--------|---------|---------------------|
| id     | string  | Unique album ID (primary key) |
| title  | string  | Album title         |
| artist | string  | Artist name         |
| price  | float64 | Price of the album  |

## Getting Started

### Prerequisites

- Go 1.18+ installed
- PostgreSQL database running and accessible
- `go-pg/pg` and `joho/godotenv` Go packages installed
- `.env` file in project root with the following variables:
        DB_HOST=localhost
        DB_PORT=5432
        DB_USER=your_pg_username
        DB_PASSWORD=your_pg_password
        DB_NAME=your_database_name

### Database Setup

Before running the server, create the `albums` table in your PostgreSQL database:

            CREATE TABLE albums (

                id VARCHAR PRIMARY KEY,
                title VARCHAR NOT NULL,
                artist VARCHAR NOT NULL,
                price NUMERIC(10,2) NOT NULL

            );

### Running the Server

    BY writing "go run main.go" start runing the server

### Create a New Album

curl -X POST -H "Content-Type: application/json" -d '{"id":"your_id","title":"your_title","artist":"artist_name","price":your_price}' http://localhost:8080/albums

### Get All Albums

  write this by open other therminal git bash " curl http://localhost:8080/albums "

### Get Album by ID

curl http://localhost:8080/albums/<your_id>

### Delete Album by ID

curl -X DELETE http://localhost:8080/albums/<your_id>