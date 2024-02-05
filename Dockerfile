# syntax=docker/dockerfile:1

FROM golang:1.21 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app_main .

FROM alpine:3.19 AS release

RUN apk add --no-cache ffmpeg

WORKDIR /app
COPY --from=build /app/app_main /app/app_main

EXPOSE 3000
ENTRYPOINT ["./app_main"]
