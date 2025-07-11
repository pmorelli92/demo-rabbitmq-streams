FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
COPY api/ ./api/
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o order ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/order .

CMD ["./order"]