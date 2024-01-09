FROM golang:1.20-bullseye AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY *.go ./
RUN go build -o /frontend

# Create a new release build stage
FROM gcr.io/distroless/base-debian11

# Set the working directory to the root directory path
WORKDIR /

# Copy over the binary built from the previous stage
COPY --from=builder /frontend /frontend

EXPOSE 8443

ENTRYPOINT ["/frontend"]
