package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/takumakume/sbomreport-to-dependencytrack/uploader"
)

type Runner interface {
	Run() error
}

type Server struct {
	uploader uploader.Uploader
	port     int
}

func New(u uploader.Uploader, port int) *Server {
	return &Server{
		uploader: u,
		port:     port,
	}
}

func (s *Server) Run() error {
	ctx := context.Background()
	http.HandleFunc("/", uploadFunc(ctx, s.uploader))
	http.HandleFunc("/healthz", healthzFunc())

	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}
	return nil
}

func healthzFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DEBUG: server.healthzFunc: ok")
		if _, err := w.Write([]byte("ok")); err != nil {
			log.Printf("DEBUG: server.healthzFunc: write failed: %v", err)
		}
	}
}

func uploadFunc(ctx context.Context, u uploader.Uploader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Printf("ERROR: server.uploadFunc: method not allowed: %s\n", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.Body == nil {
			log.Println("ERROR: server.uploadFunc: request body is empty")
			http.Error(w, "request body is empty", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("ERROR: server.uploadFunc: request body read failed: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := u.Run(ctx, body); err != nil {
			log.Printf("ERROR: server.uploadFunc: upload failed: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write([]byte("ok")); err != nil {
			log.Printf("DEBUG: server.upload: write failed: %v", err)
		}
	}
}
