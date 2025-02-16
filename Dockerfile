FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN  go build -o coinshop coin/coin.go

CMD ["/app/coinshop"]
