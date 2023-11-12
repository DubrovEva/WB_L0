package http_server

import (
	"context"
	"fmt"
	"github.com/DubrovEva/WB_L0/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var dbInstance *db.Database

type HttpServer struct {
	listener net.Listener
	server   http.Server
}

func Start(db *db.Database, URL string) (*HttpServer, error) {
	var hs HttpServer

	listener, err := net.Listen("tcp", URL)
	if err != nil {
		return nil, fmt.Errorf("Error occurred: %s", err.Error())
	}
	hs.listener = listener

	httpHandler := newHandler(db)
	hs.server = http.Server{
		Handler: httpHandler,
	}
	go func() {
		err := hs.server.Serve(listener)
		if err != nil {
			log.Fatalf("Error occurred: %s", err.Error())
		}
	}()

	log.Printf("Service run by URL: %s\n", URL)

	return &hs, nil
}

func (hs *HttpServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := hs.server.Shutdown(ctx); err != nil {
		log.Printf("Could not shut down publisher correctly: %v\n", err)
		os.Exit(1)
	}
	return nil
}

func newHandler(db *db.Database) http.Handler {

	router := chi.NewRouter()
	dbInstance = db
	router.MethodNotAllowed(methodNotAllowedHandler)
	router.NotFound(notFoundHandler)
	router.Get("/get", getOrder)
	router.Get("/search", search)
	return router
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(405)
	render.Render(w, r, ErrMethodNotAllowed)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(404)
	render.Render(w, r, ErrNotFound)
}
