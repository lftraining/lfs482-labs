package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			renderHTML(w)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/makeRequest", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			makeHTTPRequest(w)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server started on :8443")

	// Create a Server instance to listen on port 8443 with the TLS config
	server := &http.Server{
		Addr: ":8443",
	}

	log.Fatal(server.ListenAndServeTLS("/certs/frontend.pem", "/certs/frontend-key.pem"))
}

func renderHTML(w http.ResponseWriter) {
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Frontend Service Example</title>
	</head>
	<body>
		<h1>Click the button to make an HTTPS request to our backend</h1>
		<button onclick="makeRequest()">Make Request</button>
		<p id="response"></p>
		<script>
			function makeRequest() {
				var xhttp = new XMLHttpRequest();
				xhttp.onreadystatechange = function() {
					if (this.readyState == 4 && this.status == 200) {
						document.getElementById("response").innerHTML = this.responseText;
					}
				};
				xhttp.open("GET", "/makeRequest", true);
				xhttp.send();
			}
		</script>
	</body>
	</html>
	`

	tmplParsed, err := template.New("index").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmplParsed.Execute(w, nil)
}

func makeHTTPRequest(w http.ResponseWriter) {
	targetURL := "https://backend/hello"

	// Create a CA certificate pool and add the server cert chain pem file to it
	caCert, err := os.ReadFile("/certs/backend-cert-chain.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Read the key pair to create certificate
	cert, err := tls.LoadX509KeyPair("/certs/frontend.pem", "/certs/frontend-key.pem")
	if err != nil {
		log.Fatal(err)
	}

	// Create a HTTPS client and supply the created CA pool
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	resp, err := client.Get(targetURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// read response body
	body, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		fmt.Println(error)
	}
	// close response body
	resp.Body.Close()

	// print response body
	w.Write([]byte(body))
}
