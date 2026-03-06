# build stage
FROM golang:alpine AS builder
RUN apk add --no-cache

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o team_service ./main.go

# stage 2
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/team_service .
RUN chmod +x /app/team_service
ENTRYPOINT ["./team_service"]
