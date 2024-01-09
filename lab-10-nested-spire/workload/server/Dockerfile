FROM golang:1.21 AS builder

RUN git clone https://github.com/spiffe/spiffe-helper
RUN cd spiffe-helper && \
    git checkout v0.7.0 && \
    make

FROM python:3.9-slim

USER root

# Set up spiffe-helper by copying the binary from build-stage and config from the local directory
COPY --from=builder /go/spiffe-helper/spiffe-helper /opt/spire/bin/spiffe-helper
COPY helper.conf /opt/spire/config/helper.conf

WORKDIR /app

COPY main.py .

EXPOSE 8443

# Start the spiffe-helper at entry
ENTRYPOINT ["/opt/spire/bin/spiffe-helper"]
CMD ["-config", "/opt/spire/config/helper.conf"]