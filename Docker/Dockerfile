FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod tidy

COPY . .

RUN go build -o luthibox .

EXPOSE 8080

CMD ["./luthibox"]