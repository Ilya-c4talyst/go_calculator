FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o myapp

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/cmd/myapp ./myapp

COPY index.html ./index.html
COPY schema.jpeg ./schema.jpeg

COPY .env ./.env

EXPOSE 8080

CMD ["./myapp"]