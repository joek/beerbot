package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/joek/robotwebhandlers/ws"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/raspi"
	"github.com/joek/beerbot/gobot/beerbot"
	"github.com/joek/picoborgrev"
	"github.com/joek/robotwebhandlers/webcam"
)

func main() {
	var addr = flag.String("addr", ":8080", "http service address")
	var webcamHost = flag.String("webcamHost", "localhost", "Host of webcam image.")
	var webcamPort = flag.Uint("webcamPort", 8080, "Port of webcam image.")

	flag.Parse()

	com := make(chan *ws.BotCommand)
	h := ws.NewHub(com)
	go h.Run()

	gbot := gobot.NewGobot()
	gbot.AutoStop = false

	r := raspi.NewRaspiAdaptor("raspi")
	motorA := picoborgrev.NewDriver(r, "motorA", 10)
	motorB := picoborgrev.NewDriver(r, "motorB", 11)
	beer := beerbot.NewBeerBotDriver(r, "rev", motorA, motorB)

	work := func() {

		go func() {
			for c := range com {
				// TODO: Input validation
				if c.Motor != nil {
					beer.SetMotorLeft(c.Motor.Left)
					beer.SetMotorRight(c.Motor.Right)
				} else if c.Event == "Disconnect" {
					beer.Halt()
				}
			}
		}()
	}

	robot := gobot.NewRobot("beerbot",
		[]gobot.Connection{r},
		[]gobot.Device{beer},
		work,
	)

	gbot.AddRobot(robot)

	go gbot.Start()
	defer gbot.Stop()

	webcamURL := fmt.Sprintf("%s:%d", *webcamHost, *webcamPort)
	wh := webcam.NewHandler(
		webcamURL,
	)

	http.HandleFunc("/webcam", func(w http.ResponseWriter, r *http.Request) { wh.Handle(w, r) })
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { h.ServeWs(w, r) })
	http.Handle("/", http.FileServer(http.Dir("./assets")))

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	// - Sensor output (broadcast)
}
