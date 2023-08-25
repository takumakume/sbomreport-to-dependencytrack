package server

import (
	"context"
	"io"
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
	return nil
}

func healthzFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
		w.WriteHeader(http.StatusOK)
	}
}

func uploadFunc(ctx context.Context, u uploader.Uploader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.Body == nil {
			http.Error(w, "request body is empty", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := u.Run(ctx, body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("ok"))
		w.WriteHeader(http.StatusOK)
	}
}
