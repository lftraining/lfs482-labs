import socket
import ssl
import json
import os

SVID_CERT = os.getenv('SVID_CERT', '/var/run/secrets/svids/svid.0.pem')
SVID_KEY = os.getenv('SVID_KEY', '/var/run/secrets/svids/svid.0.key')
TRUST_BUNDLE = os.getenv('TRUST_BUNDLE', '/var/run/secrets/svids/federated_bundle.0.0.pem')
PORT = os.getenv('PORT', '443')

def fetch():
    context = ssl.create_default_context(ssl.Purpose.SERVER_AUTH)
    context.load_cert_chain(certfile=SVID_CERT, keyfile=SVID_KEY)
    context.check_hostname = False

    context.verify_mode = ssl.CERT_REQUIRED
    context.load_verify_locations(cafile=TRUST_BUNDLE)

    client_socket = socket.socket()
    client_socket.connect(("server", int(PORT)))

    conn = context.wrap_socket(client_socket)

    response = b""
    while True:
        data = conn.recv(4096)
        if not data:
            break
        response += data

    ship_manifest = json.loads(response.decode("utf-8"))
    print("Received ship manifest:", ship_manifest)

    conn.close()

if __name__ == "__main__":
    fetch()