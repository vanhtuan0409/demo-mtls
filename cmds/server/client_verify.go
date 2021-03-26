package main

import (
	"crypto/x509"
	"errors"
)

var (
	allowedClient = map[string]bool{
		"client.yolo": true,
	}

	ErrorInvalidClient = errors.New("Invalid client")
)

// verifyClient Do your client verification here
func verifyClient(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	for _, chain := range verifiedChains {
		if len(chain) < 1 {
			continue
		}
		cert := chain[0] //check leaf cert only
		if !allowedClient[cert.Subject.CommonName] {
			continue
		} else {
			return nil
		}
	}
	return ErrorInvalidClient
}
