FROM golang:1.15

RUN mkdir -p /app
ADD . /app
WORKDIR /app

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o app ./cmd/main.go

CMD [ "./app" ]    