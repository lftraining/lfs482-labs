ARG GOLANG_BUILDER_IMAGE
ARG GOLANG_IMAGE
FROM $GOLANG_BUILDER_IMAGE AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY *.go ./
RUN go build -o /frontend

# Create a new release build stage
FROM $GOLANG_IMAGE

# Set the working directory to the root directory path
WORKDIR /

# Copy over the binary built from the previous stage
COPY --from=builder /frontend /frontend

EXPOSE 8443

ENTRYPOINT ["/frontend"]
