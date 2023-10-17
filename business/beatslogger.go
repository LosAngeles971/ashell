package business

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// A log entry can refer to the execution of a test suite or a test case
type CommandEvent struct {
	Timestamp string `json:"@timestamp"`
	Command   string `json:"command"`
	Output    string `json:"output,omitempty"`
}

type beatsLogger struct {
	logFilename      string
	logFile          *os.File
	buffer           *bytes.Buffer
}

type beatsLoggerOption func(*beatsLogger)

func newBeatsLogger(filename string, opts ...beatsLoggerOption) *beatsLogger {
	l := &beatsLogger{
		logFilename: filename,
		buffer:         &bytes.Buffer{},
	}
	var err error
	l.logFile, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0775)
	if err != nil {
		panic(err)
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// logEvent logs a command event to a file
func (l *beatsLogger) logEvent(line []byte) {
	event := CommandEvent{
		Timestamp: time.Now().Format(time.RFC3339),
		Command:   "",
		Output:    base64.StdEncoding.EncodeToString(line),
	}
	j, err := json.Marshal(event)
	if err != nil {
		log.Error(err)
		return
	}
	f, err := os.OpenFile(l.logFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0775)
	if err != nil {
		log.Error(err)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(string(j) + "\n"); err != nil {
		log.Error(err)
	}
}

func (l *beatsLogger) Write(p []byte) (n int, err error) {
	n = len(p)
	l.buffer.Write(p)
	if bytes.Contains(p, []byte{'\n'}) {
		for _, line := range bytes.Split(l.buffer.Bytes(), []byte{'\n'}) {
			if len(line) == 0 {
				continue
			}
			l.logEvent(line)
		}
		l.buffer.Reset()
	}
	return n, nil
}

func (l *beatsLogger) Close() error {
	return l.logFile.Close()
}
