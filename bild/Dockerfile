FROM golang:1.15 as modules

ADD go.mod go.sum . /
RUN go mod download

FROM golang:1.15 as builder

COPY --from=modules /go/pkg /go/pkg

RUN mkdir -p /project
ADD . /project
WORKDIR /project

RUN useradd -u 10001 project

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
go build -o project ./cmd/main.go

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
USER project

COPY --from=builder /project /project


CMD [ "./project" ]    
