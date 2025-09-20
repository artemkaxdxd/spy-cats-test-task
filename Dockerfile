FROM golang:1.25.1-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o spy-cats-test-task ./

FROM alpine:3.17
RUN apk add --no-cache ca-certificates

WORKDIR /root/
COPY --from=builder /app/spy-cats-test-task .

EXPOSE 8080

ENTRYPOINT ["./spy-cats-test-task"]
