# # FROM golang:1.19 as modules

# # ADD go.mod go.sum ./
# # RUN go mod download

# FROM golang:1.18

# # COPY --from=modules /go/pkg /go/pkg

# RUN mkdir -p /project
# ADD . /project
# WORKDIR /project

# RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
# go build -o project ./cmd/main.go

# FROM scratch

# COPY --from=builder /project /project


# CMD [ "./project" ]    
FROM golang:1.18

WORKDIR /usr/src/app

COPY go.mod go.sum  ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/chatty cmd/main.go

CMD ["chatty"]