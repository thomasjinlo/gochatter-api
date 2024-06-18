package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/redis/go-redis/v9"
	"github.com/thomasjinlo/gochatter-api/internal/handlers"
	"github.com/thomasjinlo/gochatter-api/internal/ws"
)

func main() {
	log.Print("[gochatter-api] Starting up API Server")
	rootPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[gochatter-api] root path %s", rootPath)
	rc := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	certPool := x509.NewCertPool()
	rootCA, err := os.ReadFile(filepath.Join(rootPath, ".credentials", "root-ca.crt"))
	if err != nil {
		log.Printf("[gochatter-api] error reading ws cert: %v", err)
	}
	certPool.AppendCertsFromPEM(rootCA)
	cert, err := tls.LoadX509KeyPair(
		filepath.Join(rootPath, ".credentials", "cert.pem"),
		filepath.Join(rootPath, ".credentials", "key.pem"),
	)
	tlsConfig := &tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{cert},
	}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	hc := &http.Client{Transport: transport}
	wsClient := ws.NewClient(rc, hc)
	log.Fatal(http.ListenAndServeTLS(
		":8443",
		filepath.Join(rootPath, ".credentials", "cert.pem"),
		filepath.Join(rootPath, ".credentials", "key.pem"),
		handlers.SetupRoutes(wsClient)))

}
