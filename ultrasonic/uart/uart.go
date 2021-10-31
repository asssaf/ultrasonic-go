package uart

import (
	"encoding/binary"
	"time"

	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/uart"

	"github.com/asssaf/ultrasonic-go/ultrasonic"
)

const (
	cmdDistance byte = 0x55
	cmdTemp     byte = 0x50
)

type dev struct {
	c conn.Conn
}

var _ ultrasonic.Dev = &dev{}

// create a new sensor connected by a serial interface, such as US100 with the jumper on
func NewUart(port uart.Port) (*dev, error) {
	conn, err := port.Connect(9600*physic.Hertz, uart.One, uart.NoParity, uart.NoFlow, 8)
	if err != nil {
		return nil, err
	}

	dev := &dev{
		c: conn,
	}

	return dev, nil
}

func (d *dev) Sense(values *ultrasonic.SensorValues) error {
	distance, err := d.senseDistance()
	if err != nil {
		return err
	}

	temp, err := d.senseTemperature()
	if err != nil {
		return err
	}

	values.Distance = distance
	values.Temperature = temp

	return nil
}

func (d *dev) senseDistance() (physic.Distance, error) {
	// r/w Tx doesn't work well, need to wait 100 milliseconds between write and read
	if err := d.c.Tx([]byte{cmdDistance}, []byte{}); err != nil {
		return 0, err
	}

	time.Sleep(100 * time.Millisecond)

	read := make([]byte, 2)
	if err := d.c.Tx([]byte{}, read); err != nil {
		return 0, err
	}

	time.Sleep(100 * time.Millisecond)

	distanceRaw := binary.BigEndian.Uint16(read)
	distance := physic.Distance(distanceRaw) * physic.MilliMetre * 10
	// if cap > 4095 {
	// 	return 0, errors.New(fmt.Sprintf("bad sample: %d", cap))
	// }

	return distance, nil
}

func (d *dev) senseTemperature() (physic.Temperature, error) {
	// r/w Tx doesn't work well, need to wait 100 milliseconds between write and read
	if err := d.c.Tx([]byte{cmdTemp}, []byte{}); err != nil {
		return 0, err
	}

	time.Sleep(100 * time.Millisecond)

	read := make([]byte, 1)
	if err := d.c.Tx([]byte{}, read); err != nil {
		return 0, err
	}

	time.Sleep(100 * time.Millisecond)

	tempRaw := read[0]
	temp := physic.ZeroCelsius + physic.Temperature(tempRaw)*physic.Celsius

	return temp, nil
}

func (d *dev) Halt() error {
	return nil
}
