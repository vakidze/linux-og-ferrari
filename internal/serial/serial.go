package serial

import (
	"bufio"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/tarm/serial"
)

// ListPorts returns a list of likely ports (user can edit)
func ListPorts() ([]string, error) {
	ports := []string{"/dev/ttyUSB0", "/dev/ttyUSB1", "/dev/ttyACM0", "/dev/ttyS0"}
	return ports, nil
}

// Open opens a serial port and returns its handle
func Open(port string, baud int) (*serial.Port, error) {
	cfg := &serial.Config{
		Name: port,
		Baud: baud,
		ReadTimeout: time.Second * 2,
	}
	return serial.OpenPort(cfg)
}

// ReadLinesToChan opens the port and streams cleaned lines into a channel.
func ReadLinesToChan(port string, baud int, filterFn func(string) bool) (<-chan string, error) {
	if port == "" {
		return nil, errors.New("empty port")
	}

	s, err := Open(port, baud)
	if err != nil {
		return nil, err
	}

	out := make(chan string, 200)
	go func() {
		defer s.Close()
		defer close(out)

		reader := bufio.NewReader(s)
		ansi := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return // stop on read error
			}

			line = ansi.ReplaceAllString(line, "") // strip ANSI colors
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if filterFn != nil && !filterFn(line) {
				continue
			}

			out <- line
		}
	}()
	return out, nil
}
