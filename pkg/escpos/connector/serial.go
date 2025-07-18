package connector

import (
	"fmt"
	"io"
	"time"

	"go.bug.st/serial"
)

// SerialConfig holds all configuration for a serial port connection
type SerialConfig struct {
	BaudRate int             // Default: 19200
	DataBits int             // Default: 8
	StopBits serial.StopBits // Default: 1
	Parity   serial.Parity   // Default: none
	Timeout  time.Duration   // Default: 1s
}

// DefaultSerialConfig returns a config with sane defaults for ESC/POS printers
func DefaultSerialConfig() *SerialConfig {
	return &SerialConfig{
		BaudRate: 19200,
		DataBits: 8,
		StopBits: serial.OneStopBit,
		Parity:   serial.NoParity,
		Timeout:  1 * time.Second,
	}
}

// SerialConnector implements serial port communication for ESC/POS printers
// across Windows, Linux and macOS
type SerialConnector struct {
	portName string
	config   *SerialConfig
	port     serial.Port
	isOpen   bool
}

func NewSerialConnector(portName string, config *SerialConfig) (*SerialConnector, error) {
	if portName == "" {
		return nil, fmt.Errorf("port name cannot be empty")
	}

	if config == nil {
		config = DefaultSerialConfig()
	}

	connector := &SerialConnector{
		portName: portName,
		config:   config,
		isOpen:   false,
	}

	// Open the serial port
	if err := connector.open(); err != nil {
		return nil, err
	}

	return connector, nil
}

func (c *SerialConnector) open() error {
	mode := &serial.Mode{
		BaudRate: c.config.BaudRate,
		DataBits: c.config.DataBits,
		StopBits: c.config.StopBits,
		Parity:   c.config.Parity,
	}

	port, err := serial.Open(c.portName, mode)
	if err != nil {
		return fmt.Errorf("failed to open serial port %s: %w", c.portName, err)
	}

	// Set read timeout
	if err := port.SetReadTimeout(c.config.Timeout); err != nil {
		port.Close()
		return fmt.Errorf("failed to set read timeout: %w", err)
	}

	c.port = port
	c.isOpen = true
	return nil
}

// Write sends data to the serial port
// Implements the io.Writer interface
func (c *SerialConnector) Write(data []byte) (int, error) {
	if !c.isOpen || c.port == nil {
		return 0, fmt.Errorf("serial port is not open")
	}

	return c.port.Write(data)
}

// Close closes the serial port connection
func (c *SerialConnector) Close() error {
	if c.isOpen && c.port != nil {
		c.isOpen = false
		return c.port.Close()
	}
	return nil
}

// ReadByte reads a single byte from the serial port
// Useful for printers that send status information back
func (c *SerialConnector) Read(buf []byte) (int, error) {
	if !c.isOpen || c.port == nil {
		return 0, fmt.Errorf("serial port is not open")
	}

	_, err := c.port.Read(buf)
	if err != nil {
		if err == io.EOF {
			return 0, err
		}
		return 0, fmt.Errorf("failed to read from serial port: %w", err)
	}

	return int(buf[0]), nil
}

// Flush ensures all data is written to the serial port
func (c *SerialConnector) Flush() error {
	// Implemented only if the underlying library supports it
	// Currently go.bug.st/serial doesn't have a direct flush method,
	// but for completeness we could add this method
	return nil
}

// ListPorts returns a list of available serial ports on the system
func ListPorts() ([]string, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, fmt.Errorf("failed to list serial ports: %w", err)
	}
	return ports, nil
}

func IsPortInUse(portName string) bool {
	mode := &serial.Mode{
		BaudRate: 9600, // Configuración básica
	}
	port, err := serial.Open(portName, mode)
	if err != nil {
		return true // El puerto está en uso o no disponible
	}
	defer port.Close()
	return false // El puerto está disponible
}
