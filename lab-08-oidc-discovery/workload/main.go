package main

import (
	"context"
	"encoding/json"
	"github.com/avast/retry-go"
	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/spiffe/go-spiffe/v2/svid/jwtsvid"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"io"
	"log"
	"net/http"
)

const (
	audience = "lab-09"
)

func main() {
	ctx := context.Background()

	// Get a client for the workload api
	workloadApiClient, err := workloadapi.New(ctx)
	if err != nil {
		log.Fatalf("Failed to create workload api client: %v", err)
	}

	// Get a jwt svid, setting the desired audience
	// Retry to allow the spiffe controller manager to see the pod, and register the workload
	var svid *jwtsvid.SVID
	err = retry.Do(
		func() error {
			svid, err = workloadApiClient.FetchJWTSVID(ctx, jwtsvid.Params{Audience: audience})
			if err != nil {
				return err
			}

			return nil
		},
		retry.Attempts(25),
		retry.Context(ctx),
	)
	if err != nil {
		_ = workloadApiClient.Close()
		log.Fatalf("Failed to retrieve a jwt svid: %v", err)
	}

	// Close the workload api client
	_ = workloadApiClient.Close()

	// Convert the jwt svid into a string and print
	svidStr := svid.Marshal()
	printSvidString(svidStr)

	// Parse the string into a jwt.JSONWebToken
	parsedToken, err := jwt.ParseSigned(svidStr)
	if err != nil {
		log.Fatal(err)
	}

	// Get the headers and claims from the parsed token and print
	headers := parsedToken.Headers[0]
	claims := &jwt.Claims{}
	parsedToken.UnsafeClaimsWithoutVerification(claims)
	printParsedToken(headers, claims)

	// Get the oidc discovery document
	resp, err := http.Get("http://spire-spiffe-oidc-discovery-provider.spire/.well-known/openid-configuration")
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("OIDC Discovery Document from http://spire-spiffe-oidc-discovery-provider.spire/.well-known/openid-configuration: \n%s\n\n", string(bytes))

	// And we're only interested in the jwks uri
	var discoveryDocument DiscoveryDocument
	err = json.Unmarshal(bytes, &discoveryDocument)
	if err != nil {
		log.Fatal(err)
	}

	// Get the jwks
	resp, err = http.Get(discoveryDocument.JWKSUri)
	if err != nil {
		log.Fatal(err)
	}
	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("JSON Web Key Set: %s\n\n", string(bytes))

	// Unmarshall the response into a jose.JSONWebKeySet
	var keySet jose.JSONWebKeySet
	err = json.Unmarshal(bytes, &keySet)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the SVID string into a jose.JSONWebSignature
	jws, err := jose.ParseSigned(svidStr)
	if err != nil {
		log.Fatal(err)
	}

	// Get the key that signed the svid
	key := keySet.Key(parsedToken.Headers[0].KeyID)[0]

	// Verify the signature on the jwt svid
	verified, err := jws.Verify(key)
	if err != nil {
		log.Fatal(err)
	}

	// Print the verified claims
	log.Printf("Verified claims:\n%s\n\n", string(verified))
}

type DiscoveryDocument struct {
	JWKSUri string `json:"jwks_uri"`
}

func printParsedToken(headers jose.Header, claims *jwt.Claims) {
	log.Printf("Parsed Token\nHeaders: \nalg: %s\nkid: %s\nClaims: \niss: %s\nsub: %s\naud: %s\niat: %s\nexp: %s\n\n", headers.Algorithm, headers.KeyID, claims.Issuer, claims.Subject, claims.Audience, claims.IssuedAt.Time(), claims.Expiry.Time())
}

func printSvidString(svidStr string) {
	log.Printf("JWT SVID\n%s\n\n", svidStr)
}
