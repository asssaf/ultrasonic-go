package cmd

import (
	"flag"
	"fmt"
	"log"

	"periph.io/x/conn/v3/uart/uartreg"
	"periph.io/x/host/v3"

	"github.com/asssaf/ultrasonic-go/ultrasonic"
)

type ReadCommand struct {
	fs   *flag.FlagSet
	addr int
}

func NewReadCommand() *ReadCommand {
	c := &ReadCommand{
		fs: flag.NewFlagSet("read", flag.ExitOnError),
	}

	return c
}

func (c *ReadCommand) Name() string {
	return c.fs.Name()
}

func (c *ReadCommand) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	flag.Usage = c.fs.Usage

	return nil
}

func (c *ReadCommand) Execute() error {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	uartPort, err := uartreg.Open("")
	if err != nil {
		log.Fatal(err)
	}

	dev, err := ultrasonic.NewUart(uartPort)
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Halt()

	values := ultrasonic.SensorValues{}
	if err := dev.Sense(&values); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Distance: %d\n", values.Distance)
	fmt.Printf("Temperature: %s\n", values.Temperature)

	return nil
}
