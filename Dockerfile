# Stage 1: Build 
FROM golang:1.17 AS build

WORKDIR /app
COPY . .

RUN go build -o myapp

# Stage 2: Minimal runtime image
FROM alpine:latest

WORKDIR /app

COPY --from=build /app/myapp .

EXPOSE 8080
EXPOSE 9090

CMD ["./myapp"]
