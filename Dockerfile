# build stage
FROM golang:1.21-alpine3.18 AS build

# set working directory
WORKDIR /app

# copy source code
COPY . .

# install golang-migrate
RUN wget -O - https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz

# install dependencies
RUN go mod download

# build binary
RUN go build -o ./bin/gopos ./cmd/http/main.go

# final stage
FROM alpine:3.18 AS final

# set working directory
WORKDIR /app

# copy golang-migrate binary and migration files
COPY --from=build /app/migrate /usr/local/bin/migrate
COPY --from=build /app/internal/adapter/repository/postgres/migrations /app/migrations

# copy binary
COPY --from=build /app/bin/gopos ./

# copy entrypoint and make it executable
COPY --from=build /app/entrypoint.sh ./
RUN chmod +x /app/entrypoint.sh

EXPOSE 8080

ENTRYPOINT [ "/app/entrypoint.sh" ]
