FROM golang:1.22-alpine

WORKDIR /

RUN apk add --no-cache bash curl git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server

EXPOSE 8080

CMD ["./server"]
