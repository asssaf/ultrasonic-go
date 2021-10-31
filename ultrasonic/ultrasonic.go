package ultrasonic

import (
	"periph.io/x/conn/v3/physic"
)

type Dev interface {
	Sense(values *SensorValues) error
	Halt() error
}

type SensorValues struct {
	Distance    physic.Distance
	Temperature physic.Temperature
}
