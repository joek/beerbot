package beerbot_test

import (
	. "github.com/joek/beerbot/gobot/beerbot"
	"github.com/joek/picoborgrev/revtesthelpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Beerbot", func() {
	var motorA *revtesthelpers.FakeRevDriver
	var motorB *revtesthelpers.FakeRevDriver

	BeforeEach(func() {
		motorA = revtesthelpers.NewFakeRevDriver()
		motorB = revtesthelpers.NewFakeRevDriver()
	})

	It("Creates a new BeerBotDriver instance", func() {
		var d *BeerBotDriver
		d = NewBeerBotDriver(revtesthelpers.NewI2cTestAdaptor("adaptor"), "Test", motorA, motorB)
		Ω(d).Should(BeAssignableToTypeOf(&BeerBotDriver{}))
	})

	It("Is starting the robot", func() {
		m1 := false
		m2 := false
		epo1 := false
		epo2 := false
		motorA.StartImpl = func() []error {
			m1 = true
			return nil
		}

		motorB.StartImpl = func() []error {
			m2 = true
			return nil
		}

		motorA.ResetEPOImpl = func() error {
			epo1 = true
			return nil
		}

		motorB.ResetEPOImpl = func() error {
			epo2 = true
			return nil
		}

		d := NewBeerBotDriver(revtesthelpers.NewI2cTestAdaptor("adaptor"), "Test", motorA, motorB)
		d.Start()

		Ω(m1).Should(BeTrue())
		Ω(m2).Should(BeTrue())
		Ω(epo1).Should(BeTrue())
		Ω(epo2).Should(BeTrue())

	})

	It("Is stopping the robot", func() {
		stop1 := false
		stop2 := false
		motorB.HaltImpl = func() []error {
			stop1 = true
			return nil
		}

		motorA.HaltImpl = func() []error {
			stop2 = true
			return nil
		}

		d := NewBeerBotDriver(revtesthelpers.NewI2cTestAdaptor("adaptor"), "Test", motorA, motorB)
		d.Halt()

		Ω(stop1).Should(BeTrue())
		Ω(stop2).Should(BeTrue())
	})

	It("Is returning name", func() {

		d := NewBeerBotDriver(revtesthelpers.NewI2cTestAdaptor("adaptor"), "Test", motorA, motorB)
		d.Halt()

		Ω(d.Name()).Should(Equal("Test"))
	})

	It("Is returning connection", func() {
		i := revtesthelpers.NewI2cTestAdaptor("adaptor")
		d := NewBeerBotDriver(i, "Test", motorA, motorB)
		d.Halt()

		Ω(d.Connection()).Should(Equal(i))
	})

	It("Is setting left Motors", func() {
		var m1 float32
		var m2 float32
		motorB.SetMotorAImpl = func(p float32) error {
			m1 = p
			return nil
		}

		motorA.SetMotorAImpl = func(p float32) error {
			m2 = p
			return nil
		}

		d := NewBeerBotDriver(revtesthelpers.NewI2cTestAdaptor("adaptor"), "Test", motorA, motorB)
		d.SetMotorLeft(0.32)

		Ω(m1).Should(Equal(float32(0.32)))
		Ω(m2).Should(Equal(float32(0.32)))
	})

	It("Is setting right Motors", func() {
		var m1 float32
		var m2 float32
		motorB.SetMotorBImpl = func(p float32) error {
			m1 = p
			return nil
		}

		motorA.SetMotorBImpl = func(p float32) error {
			m2 = p
			return nil
		}

		d := NewBeerBotDriver(revtesthelpers.NewI2cTestAdaptor("adaptor"), "Test", motorA, motorB)
		d.SetMotorRight(0.32)

		Ω(m1).Should(Equal(float32(-0.32)))
		Ω(m2).Should(Equal(float32(-0.32)))
	})
})
