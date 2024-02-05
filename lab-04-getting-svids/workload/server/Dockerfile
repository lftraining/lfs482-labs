FROM ghcr.io/spiffe/spiffe-helper:nightly as HELPER
FROM python:3.9-slim

USER root

# Set up spiffe-helper by copying the binary from build-stage and config from the local directory
COPY --from=HELPER /spiffe-helper /opt/spire/bin/spiffe-helper
COPY helper.conf /opt/spire/config/helper.conf

WORKDIR /app

COPY main.py .

EXPOSE 8443

# Start the spiffe-helper at entry
ENTRYPOINT ["/opt/spire/bin/spiffe-helper"]
CMD ["-config", "/opt/spire/config/helper.conf"]