package main

import (
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

//go:embed ca.crt
var caCert []byte

func main() {
	server := flag.String("server", "https://127.0.0.1:8080", "Server address")
	clientCert := flag.String("cert", "certs/bundle_client.crt", "Client cert")
	clientKey := flag.String("key", "certs/client.key", "Client key")
	flag.Parse()

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(*clientCert, *clientKey)
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}
	r, err := client.Get(*server)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", body)
}