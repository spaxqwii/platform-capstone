FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod .
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]
