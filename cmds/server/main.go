package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	caPath := flag.String("ca", "certs/ca.crt", "CA path")
	certPath := flag.String("cert", "certs/bundle_server.crt", "Cert path")
	keyPath := flag.String("key", "certs/server.key", "Key path")
	port := flag.Uint("port", 8080, "Server port")
	flag.Parse()

	// load ca cert
	caCert, err := ioutil.ReadFile(*caPath)
	if err != nil {
		panic(err)
	}
	// load server cert
	serverCert, err := tls.LoadX509KeyPair(*certPath, *keyPath)
	if err != nil {
		panic(err)
	}

	// caCert is used to verify sent client's cert
	// serverCert is used to proof that this is an authentic server
	// verifyClient is custom logic where you can check if the client is expired or blacklisted
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		ClientCAs:             caCertPool,
		ClientAuth:            tls.RequireAndVerifyClientCert,
		Certificates:          []tls.Certificate{serverCert},
		VerifyPeerCertificate: verifyClient, // on second thought, we should not do this. TLS should be used for authentication only
	}
	tlsConfig.BuildNameToCertificate()

	addr := fmt.Sprintf(":%d", *port)
	server := &http.Server{
		Addr:      addr,
		TLSConfig: tlsConfig,
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		// you can add your own authorization protocol here instead of using `verifyClient`
		msg := fmt.Sprintf("Hello %s", extractClientName(c))
		return c.String(http.StatusOK, msg)
	})
	if err := e.StartServer(server); err != nil {
		panic(err)
	}
}

func extractClientName(c echo.Context) string {
	if c.Request().TLS == nil {
		return ""
	}
	chains := c.Request().TLS.PeerCertificates
	if len(chains) < 1 {
		return ""
	}
	cert := chains[0]
	return cert.Subject.CommonName
}
