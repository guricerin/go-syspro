FROM golang:1.21

WORKDIR /go/src/github.com/foo
COPY main.go ./
COPY go.mod ./
COPY rootfs ./rootfs
RUN go mod tidy
RUN go build -o container main.go
