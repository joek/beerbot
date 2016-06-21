package testhelpers

import "github.com/hybridgroup/gobot"

type FakeRevDriver struct {
	name          string
	connection    gobot.Connection
	SetMotorAImpl func(float32) error
	SetMotorBImpl func(float32) error
	StartImpl     func() []error
	HaltImpl      func() []error
	ResetEPOImpl  func() error
	GetEPOImpl    func() (bool, error)
}

func NewFakeRevDriver() *FakeRevDriver {
	return &FakeRevDriver{
		name:       "FakeRevDriver",
		connection: newI2cTestAdaptor("I2CTest"),
		SetMotorAImpl: func(power float32) error {
			return nil
		},
		SetMotorBImpl: func(power float32) error {
			return nil
		},
		StartImpl: func() []error {
			return nil
		},
		HaltImpl: func() []error {
			return nil
		},
		GetEPOImpl: func() (bool, error) {
			return true, nil
		},
		ResetEPOImpl: func() error {
			return nil
		},
	}
}

func (b *FakeRevDriver) SetMotorA(power float32) error {
	return b.SetMotorAImpl(power)
}

func (b *FakeRevDriver) SetMotorB(power float32) error {
	return b.SetMotorBImpl(power)
}

func (b *FakeRevDriver) Start() []error {
	return b.StartImpl()
}

func (b *FakeRevDriver) Halt() []error {
	return b.HaltImpl()
}

func (b *FakeRevDriver) Name() string {
	return b.name
}

func (b *FakeRevDriver) Connection() gobot.Connection {
	return b.connection
}

func (b *FakeRevDriver) ResetEPO() error {
	return b.ResetEPOImpl()
}
func (b *FakeRevDriver) GetEPO() (bool, error) {
	return b.GetEPOImpl()
}

type i2cTestAdaptor struct {
	name         string
	I2cReadImpl  func(i int, l int) ([]byte, error)
	I2cWriteImpl func(int, []byte) error
	I2cStartImpl func() error
}

func (t *i2cTestAdaptor) I2cStart(int) (err error) {
	return t.I2cStartImpl()
}
func (t *i2cTestAdaptor) I2cRead(i int, l int) (data []byte, err error) {
	return t.I2cReadImpl(i, l)
}
func (t *i2cTestAdaptor) I2cWrite(i int, b []byte) (err error) {
	return t.I2cWriteImpl(i, b)
}
func (t *i2cTestAdaptor) Name() string             { return t.name }
func (t *i2cTestAdaptor) Connect() (errs []error)  { return }
func (t *i2cTestAdaptor) Finalize() (errs []error) { return }

func newI2cTestAdaptor(name string) *i2cTestAdaptor {
	return &i2cTestAdaptor{
		name: name,
		I2cReadImpl: func(i int, l int) ([]byte, error) {
			b := make([]byte, l, l)
			b[1] = 0x15
			return b, nil
		},
		I2cWriteImpl: func(i int, b []byte) error {
			return nil
		},
		I2cStartImpl: func() error {
			return nil
		},
	}
}
