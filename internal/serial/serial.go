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

// ReadLinesToChan opens the port and sends cleaned lines to the channel. The caller should stop the program to close the port.
func ReadLinesToChan(port string, baud int, filterFn func(string) bool) (<-chan string, error) {
	if port == "" {
		return nil, errors.New("empty port")
	}
	cfg := &serial.Config{Name: port, Baud: baud, ReadTimeout: time.Second * 2}
	s, err := serial.OpenPort(cfg)
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
				// stop on read error
				return
			}
			line = ansi.ReplaceAllString(line, "")
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
