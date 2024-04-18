package main

import (
	"log"
	"net/http"
	"os"

	"github.com/juncheoll/triton-torrent/handler"
	"github.com/juncheoll/triton-torrent/setting"
	"github.com/juncheoll/triton-torrent/src/corsController"
	"github.com/juncheoll/triton-torrent/src/logCtrlr"
	"github.com/urfave/negroni"
)

func startServer() {
	mux := handler.CreateHandler()
	handler := negroni.Classic()
	defer mux.Close()

	handler.Use(corsController.SetCors("*", "GET, POST, PUT, DELETE", "*", true))
	handler.UseHandler(mux)

	// HTTP Server Start
	logCtrlr.Log("HTTP server start.")
	http.ListenAndServe(":"+setting.ServerPort, handler)
}

func main() {

	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}
	setting.ModelsPath = pwd + "/../models"
	println(setting.ModelsPath)

	startServer()
}
