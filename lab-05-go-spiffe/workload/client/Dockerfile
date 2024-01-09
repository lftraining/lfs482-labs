FROM golang:1.21.1-alpine3.18 AS builder

WORKDIR build

COPY main.go .
COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go build main.go

FROM scratch

COPY --from=builder /go/build/main /main

ENTRYPOINT ["/main"]