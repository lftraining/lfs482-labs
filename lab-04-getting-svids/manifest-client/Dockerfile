ARG SPIFFE_HELPER_IMAGE
ARG PYTHON_IMAGE

FROM $SPIFFE_HELPER_IMAGE as HELPER
FROM $PYTHON_IMAGE

USER root

# Set up spiffe-helper by copying the binary from build-stage and config from the local directory
COPY --from=HELPER /spiffe-helper /opt/spire/bin/spiffe-helper

WORKDIR /app

COPY main.py .

# Start the spiffe-helper at entry
ENTRYPOINT ["/opt/spire/bin/spiffe-helper"]
CMD ["-config", "/opt/spire/config/helper.conf"]