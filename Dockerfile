# syntax=docker/dockerfile:1
FROM golang:1.17-alpine
WORKDIR /app/
ENV CGO_ENABLED=1

RUN apk add build-base

COPY . .

RUN go mod download && \
    go build -o ./borschplayground ./main.go

FROM docker:20.10.17-alpine3.16
WORKDIR /app/
COPY settings.local.json ./
COPY --from=0 /app/borschplayground ./
RUN apk --no-cache add ca-certificates && \
    ./borschplayground migrate

EXPOSE 8080
ENTRYPOINT ./borschplayground --address 0.0.0.0:$PORT
