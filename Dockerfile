FROM golang:1.22-alpine

WORKDIR /app
COPY . .

RUN go mod init tg-proxy || true
RUN go mod tidy
RUN go build -o proxy .

CMD ["./proxy"]
