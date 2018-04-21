package beerbot

import (
	"log"
	"sync"

	"github.com/joek/picoborgrev"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
)

// BeerBot driver interace
type BeerBot interface {
	Name() string
	Start() error
	Halt() error
	SetMotorLeft(float32) error
	SetMotorRight(float32) error
	SetMotorA(picoborgrev.RevDriver)
	SetMotorB(picoborgrev.RevDriver)
}

// BeerBotDriver struct
type BeerBotDriver struct {
	name       string
	connector  i2c.Connector
	connection i2c.Connection
	i2c.Config
	motorA picoborgrev.RevDriver
	motorB picoborgrev.RevDriver
	lock   sync.Mutex
}

// NewBeerBotDriver creates a new beerbot driver with specified name and i2c interface and MotorController adresses
// Params:
//		conn Connector - the Adaptor to use with this Driver
func NewBeerBotDriver(a i2c.Connector) *BeerBotDriver {
	b := &BeerBotDriver{
		name:      gobot.DefaultName("BeerBot"),
		connector: a,
		Config:    i2c.NewConfig(),
		lock:      sync.Mutex{},
	}

	if b.motorA == nil {
		b.motorA = picoborgrev.NewDriver(a, i2c.WithAddress(0x10))
	}

	if b.motorB == nil {
		b.motorB = picoborgrev.NewDriver(a, i2c.WithAddress(0x11))
	}
	return b
}

// Name is giving the robot name
func (d *BeerBotDriver) Name() string {
	return d.name
}

// SetName is setting bot name
func (d *BeerBotDriver) SetName(n string) {
	d.name = n
}

// SetMotorA is setting motorA driver
func (d *BeerBotDriver) SetMotorA(m picoborgrev.RevDriver) {
	d.motorA = m
}

// SetMotorA is setting motorA driver
func (d *BeerBotDriver) SetMotorB(m picoborgrev.RevDriver) {
	d.motorB = m
}

// Connection is returning the i2c connection
func (d *BeerBotDriver) Connection() gobot.Connection {
	return d.connection.(gobot.Connection)
}

// Start is starting the robot
func (d *BeerBotDriver) Start() (err error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	error := d.motorA.Start()
	if error != nil {
		log.Println("Could not start motor A")
		return error
	}

	error = d.motorB.Start()
	if error != nil {
		log.Println("Could not start motor B")
		return error
	}

	err = d.motorB.ResetEPO()
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
