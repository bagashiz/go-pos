# build stage
FROM golang:1.21-alpine3.18 AS build

# set working directory
WORKDIR /app

# copy source code
COPY . .

# install dependencies
RUN go mod download

# build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/gopos ./cmd/http/main.go

# final stage
FROM alpine:3.18 AS final
LABEL maintainer="bagashiz"

# set working directory
WORKDIR /app

# copy binary
COPY --from=build /app/bin/gopos ./

EXPOSE 8080

ENTRYPOINT [ "./gopos" ]
