FROM golang:1.18 AS build

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build

FROM debian:latest
COPY --from=build /build/docshtest /usr/bin

RUN adduser --disabled-password me
USER me
ENV USER=me
WORKDIR /home/me

ENTRYPOINT bash
