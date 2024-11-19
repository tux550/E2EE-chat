package main

import (
	"crypto/tls"
	"net/http"
)

func main() {
	server := NewWsServer()
	http.HandleFunc("/ws", server.connnect)

	// Load your certificates
	certFile := "../certs/server.crt"
	keyFile := "../certs/server.key"

	srv := &http.Server{
		Addr: "0.0.0.0:8765",
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	err := srv.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		panic(err)
	}

	//http.ListenAndServe("0.0.0.0:8765", nil)
}
