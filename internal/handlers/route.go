package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thomasjinlo/gochatter-api/internal/ws"
)

func SetupRoutes(wsClient *ws.Client) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("[gochatter-api] received hello world request")
		w.Write([]byte("Hello, World!"))
		w.WriteHeader(http.StatusOK)
	}))
	// r.Post("/broadcast", Broadcast(s))
	r.Post("/direct_message", DirectMessage(wsClient))

	return r
}

type DirectMessageRequest struct {
	SourceAccountId string
	TargetAccountId string
	Content         string
}

func DirectMessage(wsClient *ws.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[gochatter-api] received dm")
		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			e := fmt.Sprintf("unsupported content type: %s", contentType)
			log.Print(e)
			http.Error(w, e, http.StatusUnsupportedMediaType)
			return
		}
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var msg DirectMessageRequest
		err = json.Unmarshal(body, &msg)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = wsClient.SendDirectMessage(msg.TargetAccountId, msg.SourceAccountId, msg.Content)
		if err != nil {
			log.Printf("[gochatter-api] error from redis: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
