package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/juncheoll/triton-torrent/setting"
	"github.com/juncheoll/triton-torrent/tritonController"
	"github.com/unrolled/render"
)

var rend *render.Render = render.New()

type Handler struct {
	http.Handler
}

func CreateHandler() *Handler {
	mux := mux.NewRouter()
	handler := &Handler{
		Handler: mux,
	}

	downloaded := make(chan string)
	channel := &downloaded

	go tritonController.Seeding(setting.ModelsPath, channel)

	mux.HandleFunc("/serving", func(w http.ResponseWriter, r *http.Request) {
		handler.servingHandler(w, r, channel)
	}).Methods("POST") // Model serving API

	return handler
}
