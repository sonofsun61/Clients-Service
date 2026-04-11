FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp cmd/ClientService/main.go

FROM alpine:latest
COPY --from=builder /app/myapp .
EXPOSE 9191

CMD ["./myapp"]
