package logger

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"
)

type CSVLogger struct {
	file   *os.File
	writer *csv.Writer
	mu     sync.Mutex
	name   string
}

func NewCSVLogger() (*CSVLogger, error) {
	ts := time.Now().Format("20060102_150405")
	name := fmt.Sprintf("data_%s.csv", ts)
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return &CSVLogger{file: f, writer: csv.NewWriter(f), name: name}, nil
}

func (l *CSVLogger) Write(timestamp, line string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	_ = l.writer.Write([]string{timestamp, line})
	l.writer.Flush()
}

func (l *CSVLogger) CloseAndSave() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file != nil {
		l.file.Close()
		l.file = nil
	}
}

func (l *CSVLogger) Discard() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file == nil {
		return nil
	}
	name := l.name
	_ = l.file.Close()
	l.file = nil
	return os.Remove(name)
}
