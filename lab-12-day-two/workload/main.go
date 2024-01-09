package main

import (
	"context"
	"log"

	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const socketPath = "unix:///spiffe-workload-api/spire-agent.sock"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := workloadapi.New(ctx, workloadapi.WithAddr(socketPath))
	if err != nil {

		log.Fatalf("Unable to establish workload API client: %v", err)
	}
	defer client.Close()

	err = client.WatchX509Context(ctx, &x509Watcher{})
	if err != nil && status.Code(err) != codes.Canceled {
		log.Fatalf("Error watching X.509 context: %v", err)
	}
}

type x509Watcher struct{}

func (x509Watcher) OnX509ContextUpdate(c *workloadapi.X509Context) {
	for _, svid := range c.SVIDs {
		pem, _, err := svid.Marshal()
		if err != nil {
			log.Fatalf("Unable to marshal X509-SVID: %v", err)
		}
		log.Printf("Received X509-SVID for %q: \n%s\n", svid.ID, string(pem))
	}
}

func (x509Watcher) OnX509ContextWatchError(err error) {
	if status.Code(err) != codes.Canceled {
		log.Printf("X509 context watch error: %v", err)
	}
}
