import socket
import ssl
import json
import os

SVID_CERT = os.getenv('SVID_CERT', '/var/run/secrets/svids/svid.0.pem')
SVID_KEY = os.getenv('SVID_KEY', '/var/run/secrets/svids/svid.0.key')
TRUST_BUNDLE = os.getenv('TRUST_BUNDLE', '/var/run/secrets/svids/federated_bundle.0.0.pem')
PORT = os.getenv('PORT', '8443')

def serve():
    context = ssl.create_default_context(ssl.Purpose.CLIENT_AUTH)
    context.load_cert_chain(certfile=SVID_CERT, keyfile=SVID_KEY)

    context.verify_mode = ssl.CERT_REQUIRED
    context.load_verify_locations(cafile=TRUST_BUNDLE)

    server_socket = socket.socket()
    server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    server_socket.bind(("0.0.0.0", int(PORT)))
    server_socket.listen(1)

    print("HQ Server started. Waiting for connections...")

    while True:
        conn, addr = server_socket.accept()
        conn = context.wrap_socket(conn, server_side=True)
        print(f"Connection from: ", addr)

        ship_manifest = json.dumps({
            "ship_name": "SS Coastal Carrier",
            "departure_port": "London Gateway",
            "arrival_port": "Port Elizabeth",
            "cargo": [
                {"type": "electronics", "quantity": 1000},
                {"type": "clothing", "quantity": 2000},
                {"type": "food", "quantity": 3000},
            ]
        })

        conn.send(ship_manifest.encode("utf-8"))
        conn.close()
        print(f"Ship manifest sent, connection from {addr} closed.")

if __name__ == "__main__":
    serve()
