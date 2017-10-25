package main

import (
	"time"

	"github.com/joek/beerbot/gobot/beerbot"
	"github.com/joek/picoborgrev"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {

	r := raspi.NewAdaptor()
	motorA := picoborgrev.NewDriver(r, "motorA", 10)
	motorB := picoborgrev.NewDriver(r, "motorB", 11)
	beer := beerbot.NewBeerBotDriver(r, "rev", motorA, motorB)

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
