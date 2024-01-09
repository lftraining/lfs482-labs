FROM golang:1.21.1

LABEL authors="ric"

WORKDIR build

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /workload

CMD ["/workload"]
