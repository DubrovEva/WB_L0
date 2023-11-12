package http_server

import (
	"fmt"
	"github.com/go-chi/render"
	"html/template"
	"log"
	"net/http"
)

func logError(err error) {
	if err != nil {
		log.Printf("error: %s\n", err)
	}
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("id")
	if uid == "" {
		logError(render.Render(w, r, ErrorRenderer(fmt.Errorf("необходимо передать id счета"))))
		return
	}

	order, err := dbInstance.GetOrderFromCache(uid)
	if err != nil {
		logError(render.Render(w, r, ErrorRenderer(err)))
		return
	}

	err = render.Render(w, r, order)
	if err != nil {
		log.Printf("can't render order")
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	path := "internal/web/search.html"
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}
