package logger

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type CSVLogger struct {
	file  *os.File
	w     *csv.Writer
	mutex sync.Mutex
}

// Create logs directory & new .csv file using timestamp
func New(dir string) *CSVLogger {
	os.MkdirAll(dir, 0755)

	filename := time.Now().Format("2006-01-02_15-04-05") + ".csv"
	path := filepath.Join(dir, filename)

	f, _ := os.Create(path)
	w := csv.NewWriter(f)

	// Write header
	w.Write([]string{"timestamp_ms", "data"})
	w.Flush()

	return &CSVLogger{
		file: f,
		w:    w,
	}
}

// WriteLine writes log entry into csv safely
func (l *CSVLogger) WriteLine(data string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	ts := time.Now().UnixMilli()

	_ = l.w.Write([]string{time.UnixMilli(ts).Format(time.RFC3339Nano), data})
	l.w.Flush()
}

// Close closes the file
func (l *CSVLogger) Close() {
	l.w.Flush()
	l.file.Close()
}
