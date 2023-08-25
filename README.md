# Go POS

## **⚠️WIP⚠️**

## Description

A simple RESTful Point of Sale (POS) web service written in Go programming language. This project is a part of my learning process in understanding [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) in Go. The project is still in progress and will be updated from time to time.

It uses [Gin](https://gin-gonic.com/) as the HTTP framework and [PostgreSQL](https://www.postgresql.org/) as the database with [pgx](https://github.com/jackc/pgx/) as the driver and [Squirrel](https://github.com/Masterminds/squirrel) as the query builder.

## Getting Started

1. Ensure you have [Go](https://go.dev/dl/) 1.21 or higher and [Task](https://taskfile.dev/installation/) installed on your machine:

    ```bash
    go version && task --version
    ```

2. Install all required tools for the project:

    ```bash
    task install
    ```

3. Create a copy of the `.env.example` file and rename it to `.env`:

    ```bash
    cp .env.example .env
    ```

    Update configuration values as needed.

4. Run the PostgreSQL container:

    ```bash
    task db:up
    task db:create
    task migrate:up
    ```

    **NOTE**: the command use `podman-compose` by default. If you want to use `docker compose`, manually run `docker compose up -d` instead or replace `podman-compose` with `docker compose` in the `db:up` task inside `Taskfile.yml` file.

5. Run the project in development mode:

    ```bash
    task dev
    ```

## Documentation

For database schema documentation, see [here](https://dbdocs.io/bagashiz/Go-POS/), powered by [dbdocs.io](https://dbdocs.io/).

API documentation is still in progress.

## Contributing

Developers interested in contributing to Go POS project can refer to the [CONTRIBUTING](CONTRIBUTING.md) file for detailed guidelines and instructions on how to contribute.

## License

Go POS project is licensed under the [MIT License](LICENSE), providing an open and permissive licensing approach for further development and usage.

## Learning References

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) by Alistair Cockburn
- [Ready for changes with Hexagonal Architecture](https://netflixtechblog.com/ready-for-changes-with-hexagonal-architecture-b315ec967749) by Netflix Technology Blog
- [Hexagonal Architecture in Go](https://medium.com/@matiasvarela/hexagonal-architecture-in-go-cfd4e436faa3) by Matias Varela
