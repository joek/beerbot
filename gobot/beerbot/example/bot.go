package main

import (
	"time"

	"github.com/joek/beerbot/gobot/beerbot"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {

	r := raspi.NewAdaptor()

	beer := beerbot.NewBeerBotDriver(r)

	work := func() {
		beer.SetMotorLeft(0.5)
		beer.SetMotorRight(0.5)
		time.Sleep(5 * time.Second)
		beer.SetMotorLeft(-0.5)
		beer.SetMotorRight(-0.5)
		time.Sleep(5 * time.Second)
		beer.SetMotorLeft(0)
		beer.SetMotorRight(0)
	}

	robot := gobot.NewRobot("beerbot",
		[]gobot.Connection{r},
		[]gobot.Device{beer},
		work,
	)

	robot.Start()
}
