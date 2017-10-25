package beerbot

import (
	"sync"

	"github.com/joek/picoborgrev"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
)

var _ gobot.Driver = (*BeerBotDriver)(nil)

// BeerBot driver interace
type BeerBot interface {
	Name() string
	Connection() gobot.Connection
	Start() error
	Halt() error
	SetMotorLeft(float32) error
	SetMotorRight(float32) error
}

// BeerBotDriver struct
type BeerBotDriver struct {
	name       string
	connection i2c.I2c
	motorA     picoborgrev.RevDriver
	motorB     picoborgrev.RevDriver
	lock       sync.Mutex
}

// NewBeerBotDriver creates a new beerbot driver with specified name and i2c interface and MotorController adresses
func NewBeerBotDriver(a i2c.I2c, name string, motorA picoborgrev.RevDriver, motorB picoborgrev.RevDriver) *BeerBotDriver {
	return &BeerBotDriver{
		name:       name,
		connection: a,
		motorA:     motorA,
		motorB:     motorB,
		lock:       sync.Mutex{},
	}
}

// Name is giving the robot name
func (d *BeerBotDriver) Name() string {
	return d.name
}

// SetName is setting bot name
func (d *BeerBotDriver) SetName(n string) {
	d.name = n
}

// Connection is returning the i2c connection
func (d *BeerBotDriver) Connection() gobot.Connection {
	return d.connection
}

// Start is starting the robot
func (d *BeerBotDriver) Start() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	error := d.motorA.Start()
	if error != nil {
		return error
	}

	error = d.motorB.Start()
	if error != nil {
		return error
	}

	err := d.motorB.ResetEPO()
	if err != nil {
		return err
	}

	err = d.motorA.ResetEPO()
	if err != nil {
		return err
	}

	return nil
}

// Halt is stopping the robot
func (d *BeerBotDriver) Halt() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	error := d.motorA.Halt()
	if error != nil {
		return error
	}

	error = d.motorB.Halt()
	if error != nil {
		return error
	}

	return nil
}

// SetMotorLeft is setting motor speed of left motor
func (d *BeerBotDriver) SetMotorLeft(p float32) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	err := d.motorA.SetMotorA(p)
	if err != nil {
		return err
	}

	err = d.motorB.SetMotorA(p)
	if err != nil {
		return err
	}

	return nil
}

// SetMotorRight is setting motor speed of right motor
func (d *BeerBotDriver) SetMotorRight(p float32) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	p = p * (-1)

	err := d.motorA.SetMotorB(p)
	if err != nil {
		return err
	}

	err = d.motorB.SetMotorB(p)
	if err != nil {
		return err
	}

	return nil
}
