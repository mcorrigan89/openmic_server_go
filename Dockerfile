# syntax=docker/dockerfile:1

# Build
FROM golang:1.22
WORKDIR /usr/src/app

RUN apt update && apt install --no-install-recommends libvips libvips-dev pkg-config -y
COPY . .
RUN go build -o=./bin/main ./cmd

EXPOSE 8080

CMD ["./bin/main"]