FROM python:3.9-slim

USER root

WORKDIR /app

COPY main.py .

# Start the spiffe-helper at entry
CMD ["python", "main.py"]
