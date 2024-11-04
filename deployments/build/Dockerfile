
FROM golang:latest as builder
WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download -x
COPY . .
RUN go build --ldflags='-s -w -extldflags "-static"' cmd/main.go
ENTRYPOINT ["/go/src/app/main"]