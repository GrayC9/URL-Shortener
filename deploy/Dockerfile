FROM golang:alpine
COPY ../ .
RUN go build -o shortener ./cmd/main.go
CMD ["./shortener"]
