FROM golang:1.13 AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o main ./cmd/server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/main .
EXPOSE 8088
ENTRYPOINT ["./main", "-p", "8088"]