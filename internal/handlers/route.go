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

type DirectMessager interface {
	DirectMessage(accountId, content string) error
}

type Sender interface {
	Broadcaster
	DirectMessager
}

func SetupRoutes(s DirectMessager) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("[gochatter-api] received hello world request")
		w.Write([]byte("Hello, World!"))
		w.WriteHeader(http.StatusOK)
	}))
	// r.Post("/broadcast", Broadcast(s))
	r.Post("/direct_message", DirectMessage(s))

	return r
}

type BroadcastMessage struct {
	Author  string
	Content string
}

func Broadcast(s Broadcaster) http.HandlerFunc {
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
		var msg BroadcastMessage
		err = json.Unmarshal(body, &msg)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := s.Broadcast(msg.Author, msg.Content)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(res.StatusCode)
	}
}

type DirectMessageBody struct {
	AccountId string
	Content   string
}

func DirectMessage(s DirectMessager) http.HandlerFunc {
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
		var msg DirectMessageBody
		err = json.Unmarshal(body, &msg)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = s.DirectMessage(msg.AccountId, msg.Content)
		if err != nil {
			log.Printf("[gochatter-api] error from redis: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
