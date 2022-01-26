FROM golang:1.17.6-alpine3.15

WORKDIR /app
COPY /src ./src
COPY go.mod ./
COPY go.sum ./
COPY README.md ./

RUN go mod download
RUN go build -o ./src/todo_app ./src/main

RUN chmod +x ./src/todo_app
ENV GIN_MODE=release

EXPOSE 8112

ENTRYPOINT ["./src/todo_app"]