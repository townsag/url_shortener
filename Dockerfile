FROM golang:1.24 AS builder
WORKDIR /app
COPY src/go.mod src/go.sum ./
RUN go mod download
COPY src/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main

FROM alpine/curl:latest AS runner
WORKDIR /app
COPY --from=builder /app/main .
RUN adduser -D appuser
USER appuser
CMD ["./main"]