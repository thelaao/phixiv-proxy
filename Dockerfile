# syntax=docker/dockerfile:1

FROM golang:1.21-alpine AS build

RUN apk add --no-cache build-base libjpeg-turbo-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app_main .

FROM alpine:3.19 AS release

RUN apk add --no-cache ffmpeg libjpeg-turbo-dev

WORKDIR /app
COPY --from=build /app/app_main /app/app_main

EXPOSE 3000
CMD ["./app_main"]
