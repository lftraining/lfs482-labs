ARG GOLANG_BUILDER_IMAGE
ARG GOLANG_IMAGE

FROM $GOLANG_BUILDER_IMAGE AS BUILDER

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY *.go ./
RUN go build -o /workload

FROM $GOLANG_IMAGE

COPY --from=BUILDER /workload /workload

ENTRYPOINT ["/workload"]