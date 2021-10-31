package cmd

import (
	"flag"
	"fmt"
	"log"
	"time"

	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/uart/uartreg"
	"periph.io/x/host/v3"

	"github.com/asssaf/ultrasonic-go/ultrasonic"
	"github.com/asssaf/ultrasonic-go/ultrasonic/gpio"
	"github.com/asssaf/ultrasonic-go/ultrasonic/uart"
)

type ReadCommand struct {
	fs         *flag.FlagSet
	uart       string
	trigger    string
	echo       string
	continuous bool
}

func NewReadCommand() *ReadCommand {
	c := &ReadCommand{
		fs: flag.NewFlagSet("read", flag.ExitOnError),
	}

	c.fs.StringVar(&c.uart, "uart", "", "UART device (/dev/ttyS0)")
	c.fs.StringVar(&c.trigger, "trigger", "", "Trigger GPIO pin (14)")
	c.fs.StringVar(&c.echo, "echo", "", "Echo GPIO pin (15)")
	c.fs.BoolVar(&c.continuous, "continous", false, "Continous reading")

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

	var dev ultrasonic.Dev
	if c.uart != "" {
		if c.trigger != "" || c.echo != "" {
			log.Fatal("Can't use -trigger or -echo when -uart is used")
		}

		uartPort, err := uartreg.Open(c.uart)
		if err != nil {
			log.Fatal(err)
		}

		dev, err = uart.NewUart(uartPort)
		if err != nil {
			log.Fatal(err)
		}
		defer dev.Halt()
	} else {
		if c.trigger == "" || c.echo == "" {
			log.Fatal("-trigger and -echo must be set unless -uart is used")
		}

		trigger := gpioreg.ByName(c.trigger)
		if trigger == nil {
			log.Fatalf("couldn't find pin: %s", c.trigger)
		}

		echo := gpioreg.ByName(c.echo)
		if echo == nil {
			log.Fatalf("couldn't find pin: %s", c.echo)
		}

		var err error
		dev, err = gpio.NewGPIO(trigger, echo)
		if err != nil {
			log.Fatal(err)
		}
		defer dev.Halt()
	}

	for {
		values := ultrasonic.SensorValues{}
		if err := dev.Sense(&values); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Distance: %s\n", values.Distance)
		fmt.Printf("Temperature: %s\n", values.Temperature)

		if !c.continuous {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}
	return nil
}
