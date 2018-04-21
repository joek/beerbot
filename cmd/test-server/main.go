package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/joek/robotwebhandlers/webcam"
	"github.com/joek/robotwebhandlers/ws"
)

func main() {
	var addr = flag.String("addr", ":8080", "http service address")
	var webcamHost = flag.String("webcamHost", "localhost", "Host of webcam image.")
	var webcamPort = flag.Uint("webcamPort", 8080, "Port of webcam image.")
	var assetPath = flag.String("assetPath", "../server/assets", "Path to www resources.")
	var err error
	flag.Parse()

	com := make(chan *ws.BotCommand)
	h := ws.NewHub(com)
	go h.Run()

	go func() {
		for c := range com {
			log.Printf("Command: %v", c)
		}
	}()

	webcamURL := fmt.Sprintf("%s:%d", *webcamHost, *webcamPort)
	wh := webcam.NewHandler(
		webcamURL,
	)

	http.HandleFunc("/webcam", func(w http.ResponseWriter, r *http.Request) { wh.Handle(w, r) })
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { h.ServeWs(w, r) })
	http.Handle("/", http.FileServer(http.Dir(*assetPath)))

	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	// Robot WS
	// - Control input (lock to single connection)
	// - Sensor output (broadcast)
	// Webcam handler
}
