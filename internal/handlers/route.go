package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Broadcaster interface {
	Broadcast(author, content string) (*http.Response, error)
}

func SetupRoutes(b Broadcaster) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("[gochatter-api] received hello world request")
		w.Write([]byte("Hello, World!"))
		w.WriteHeader(http.StatusOK)
	}))
	r.Post("/send_message", SendMessage(b))

	return r
}

type SendMessageBody struct {
	Author string
	Content string
}

func SendMessage(b Broadcaster) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		var msg SendMessageBody
		err = json.Unmarshal(body, &msg)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := b.Broadcast(msg.Author, msg.Content)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(res.StatusCode)
	}
}
