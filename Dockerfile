FROM golang:1.23

WORKDIR /app

COPY . .

RUN go mod tidy
RUN  go build -o coinshop cmd/app/main.go

CMD [ "/app/coinshop", "--config=./cmd/app/config/config.yml" ]
