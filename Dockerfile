FROM golang:1.19

RUN mkdir /app
WORKDIR /app
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./migrations ./migrations
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

RUN go mod download

RUN GOOS=linux go build -o gophermart cmd/gophermart/*.go


EXPOSE 8080
ENTRYPOINT ["/app/gophermart"]


