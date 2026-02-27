FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN go build -o service .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/service .
EXPOSE 3001
CMD ["./service"]
