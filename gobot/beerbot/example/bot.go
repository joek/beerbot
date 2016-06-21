package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/raspi"
	"github.com/joek/beerbot/gobot/beerbot"
	"github.com/joek/picoborgrev"
)

func main() {
	gbot := gobot.NewGobot()

	r := raspi.NewRaspiAdaptor("raspi")
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

	gbot.AddRobot(robot)

	gbot.Start()
}
