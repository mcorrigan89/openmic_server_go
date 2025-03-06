# syntax=docker/dockerfile:1

# Build
FROM golang:1.23-alpine AS build
WORKDIR /app
COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o ./bin/main ./cmd

FROM gcr.io/distroless/static-debian12

COPY --from=build /bin .

EXPOSE 8080

CMD ["./bin/main"]