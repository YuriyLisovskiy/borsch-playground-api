FROM golang:1.17-alpine
WORKDIR /app/
ENV CGO_ENABLED=1

RUN apk add build-base

COPY . .

RUN go mod download
RUN go build -o ./playground_app ./cmd/playground/main.go

FROM docker:20.10.17-alpine3.16
WORKDIR /app/

COPY settings.json ./
COPY --from=0 /app/playground_app ./

RUN apk --no-cache add ca-certificates
