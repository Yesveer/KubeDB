FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache ca-certificates bash curl
WORKDIR /
COPY --from=builder /server /server
COPY --from=builder /app/scripts /scripts
EXPOSE 8080
ENTRYPOINT ["/server"]