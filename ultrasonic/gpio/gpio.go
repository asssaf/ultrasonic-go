package gpio

import (
	"fmt"
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
	// make sure output is low
	err := d.triggerPin.Out(gpio.Low)
	if err != nil {
		return 0, err
	}
	time.Sleep(60 * time.Millisecond)

	// create the trigger pulse
	err = d.triggerPin.Out(gpio.High)
	if err != nil {
		return 0, err
	}

	time.Sleep(10 * time.Microsecond)
	err = d.triggerPin.Out(gpio.Low)
	if err != nil {
		return 0, err
	}

	err = d.echoPin.In(gpio.PullDown, gpio.RisingEdge)
	if err != nil {
		return 0, err
	}

	if d.echoPin.Read() == gpio.High {
		return 0, fmt.Errorf("echo is high before timer was started")
	}

	ok := d.echoPin.WaitForEdge(time.Second)
	if !ok {
		return 0, fmt.Errorf("timeout waiting for rising edge")
	}

	start := time.Now()

	err = d.echoPin.In(gpio.PullNoChange, gpio.BothEdges)
	if err != nil {
		return 0, err
	}

	ok = d.echoPin.WaitForEdge(time.Second)
	if !ok {
		return 0, fmt.Errorf("timeout waiting for falling edge")
	}

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
