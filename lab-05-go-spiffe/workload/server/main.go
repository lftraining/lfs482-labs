package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
)

// serverAddress is the address that the server will listen on
const (
	serverAddress = "0.0.0.0:8443"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	log.Printf("SPIFFE_ENDPOINT_SOCKET: %s", os.Getenv("SPIFFE_ENDPOINT_SOCKET"))

	// Define the required SPIFFE ID of the client
	clientID := spiffeid.RequireFromString("spiffe://coastal-containers.example/workload/client")

	log.Printf("Starting server on %s", serverAddress)

	// Create a Workload API listener for the server
	listener, err := spiffetls.Listen(ctx, "tcp", serverAddress, tlsconfig.AuthorizeID(clientID))
	if err != nil {
		return fmt.Errorf("Error.. Unable to create TLS listener: %w", err)
	}
	defer listener.Close()

	log.Printf("Server listening on %s", serverAddress)
	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("Error.. Failed to accept connection: %w", err)
		}
		go handleConnection(conn)
	}
}

// handleConnection reads a request from the client and sends a response
func handleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Printf("Error reading incoming data: %v", err)
		return
	}
	log.Printf("Incoming vessel says: %q", req)

	// Send a response back to the client
	if _, err = conn.Write([]byte("Request received SS Coastal Carrier. You are cleared to dock.\n")); err != nil {
		log.Printf("Unable to send response: %v", err)
		return
	}
}
