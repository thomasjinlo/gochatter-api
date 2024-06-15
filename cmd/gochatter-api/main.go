package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/thomasjinlo/gochatter-api/internal/handlers"
	//	"github.com/thomasjinlo/gochatter-api/internal/pushserver"
	"github.com/thomasjinlo/gochatter-api/internal/send"
)

func main() {
	log.Print("[gochatter-api] Starting up API Server")

	rootPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[gochatter-api] root path %s", rootPath)

	//	psIp := os.Getenv("PUSH_SERVER_IP")
	// ps := pushserver.NewPushServer(psIp)
	s := &send.Sender{}
	log.Fatal(http.ListenAndServeTLS(
		":8443",
		filepath.Join(rootPath, ".credentials", "cert.pem"),
		filepath.Join(rootPath, ".credentials", "key.pem"),
		handlers.SetupRoutes(s)))

}
