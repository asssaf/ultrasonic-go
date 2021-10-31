package gpio

import (
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/physic"

	"github.com/asssaf/ultrasonic-go/ultrasonic"
)

type dev struct {
	triggerPin gpio.PinOut
	echoPin    gpio.PinIn
}

var _ ultrasonic.Dev = &dev{}

// create a new sensor connected by a serial interface, such as US100 with the jumper on
func NewGPIO(triggerPin gpio.PinOut, echoPin gpio.PinIn) (*dev, error) {
	dev := &dev{
		triggerPin: triggerPin,
		echoPin:    echoPin,
	}
	return dev, nil
}

func (d *dev) Sense(values *ultrasonic.SensorValues) error {
	distance, err := d.senseDistance()
	if err != nil {
		return err
	}

	values.Distance = distance

	return nil
}

func (d *dev) senseDistance() (physic.Distance, error) {
	err := d.triggerPin.Out(gpio.High)
	if err != nil {
		return 0, err
	}

	time.Sleep(10 * time.Microsecond)
	err = d.triggerPin.Out(gpio.Low)
	if err != nil {
		return 0, err
	}

	err = d.echoPin.In(gpio.PullNoChange, gpio.RisingEdge)
	if err != nil {
		return 0, err
	}

	_ = d.echoPin.WaitForEdge(-1)

	start := time.Now()

	err = d.echoPin.In(gpio.PullNoChange, gpio.FallingEdge)
	if err != nil {
		return 0, err
	}

	_ = d.echoPin.WaitForEdge(-1)

	pulseTime := time.Since(start)

	// stop edge detection
	err = d.echoPin.In(gpio.PullNoChange, gpio.NoEdge)
	if err != nil {
		return 0, err
	}

	// speed of sound is 343m/s, and we're doing round trip so divide time by two
	distanceRaw := float64(pulseTime) * 343 / 2
	// duration is in nanoseconds, so multiply by nanometer
	distance := physic.Distance(distanceRaw) * physic.NanoMetre

	return distance, nil
}

func (d *dev) Halt() error {
	return nil
}
