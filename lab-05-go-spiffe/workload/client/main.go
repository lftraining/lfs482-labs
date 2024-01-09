package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
)

// serverAddress is the address that the client will connect to
// this is set to the `server` service name to work in K8s
const (
	serverAddress = "server:443"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	log.Printf("SPIFFE_ENDPOINT_SOCKET: %s", os.Getenv("SPIFFE_ENDPOINT_SOCKET"))

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Define the required SPIFFE ID of the server
	serverID := spiffeid.RequireFromString("spiffe://coastal-containers.example/workload/server")

	log.Printf("Connecting to %s", serverAddress)
	// Dial the server to establish a connection
	conn, err := spiffetls.Dial(ctx, "tcp", serverAddress, tlsconfig.AuthorizeID(serverID))
	if err != nil {
		return fmt.Errorf("Error.. unable to establish TLS connection: %w", err)
	}
	defer conn.Close()
	log.Printf("Client connected to %s", serverAddress)

	// Send a message to the server
	fmt.Fprintf(conn, "This is SS Coastal Carrier hailing the port authority for Coastal Containers Ltd.\n")

	// Read the response from the server
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil && err != io.EOF {
		return fmt.Errorf("Error.. Unable to read server response: %w", err)
	}
	log.Printf("Port Authority says: %q", status)
	return nil
}
