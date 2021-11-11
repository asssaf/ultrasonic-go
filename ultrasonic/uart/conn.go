package uart

import (
	"fmt"
	"io"

	"github.com/jacobsa/go-serial/serial"
	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/uart"
)

var _ uart.PortCloser = &serialPort{}
var _ conn.Conn = &serialConn{}

type serialPort struct {
	name string
	c    *serialConn
}

func NewPort(name string) *serialPort {
	p := &serialPort{
		name: name,
	}

	return p
}

func (s *serialPort) String() string {
	return s.name
}

func (s *serialPort) Connect(f physic.Frequency, stopBit uart.Stop, parity uart.Parity, flow uart.Flow, bits int) (conn.Conn, error) {
	var stopBits uint
	switch stopBit {
	case uart.One:
		stopBits = 1
	case uart.Two:
		stopBits = 2
	default:
		return nil, fmt.Errorf("unsupported stop bits: %s", stopBit)
	}

	options := serial.OpenOptions{
		PortName:        s.name,
		BaudRate:        uint(f / physic.Hertz),
		DataBits:        uint(bits),
		StopBits:        stopBits,
		MinimumReadSize: 1,
	}

	// Open the port.
	h, err := serial.Open(options)
	if err != nil {
		return nil, err
	}

	c := &serialConn{
		name: s.name,
		h:    h,
	}

	s.c = c

	return c, nil
}

func (s *serialPort) Close() error {
	if s.c == nil {
		return nil
	}
	return s.c.h.Close()
}

func (p *serialPort) LimitSpeed(f physic.Frequency) error {
	return nil
}

type serialConn struct {
	name string
	h    io.ReadWriteCloser
}

func (s *serialConn) String() string {
	return s.name
}

func (s *serialConn) Tx(w, r []byte) error {
	if len(w) > 0 {
		wc, err := s.h.Write(w)
		if err != nil {
			return fmt.Errorf("writing %d err: %w", w, err)
		}
		if wc != len(w) {
			return fmt.Errorf("wrote %d bytes instead of %d", wc, len(w))
		}
	}

	if len(r) > 0 {
		rc, err := s.h.Read(r)
		if err != nil {
			return fmt.Errorf("error reading %d bytes: %w", len(r), err)
		}
		if rc != len(r) {
			return fmt.Errorf("read %d bytes instead of %d", rc, len(r))
		}
	}

	return nil
}

func (s *serialConn) Duplex() conn.Duplex {
	return conn.Full
}
