FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main

FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/main .
RUN adduser -D appuser
USER appuser
CMD ["./main"]