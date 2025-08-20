FROM node:24-alpine AS static-site-builder
WORKDIR /app
COPY ui/package.json ui/package-lock.json ./
RUN npm ci
COPY ./ui .
RUN npm run build
RUN npm prune --production

FROM golang:1.24 AS builder
WORKDIR /app
COPY api/go.mod api/go.sum ./
RUN go mod download
COPY api/ .
COPY --from=static-site-builder /app/build ./build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main

FROM alpine/curl:latest AS runner
WORKDIR /app
COPY --from=builder /app/main .
RUN adduser -D appuser
USER appuser
CMD ["./main"]