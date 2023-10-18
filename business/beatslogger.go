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
	logFilename string
	buffer      *bytes.Buffer
}

type beatsLoggerOption func(*beatsLogger)

func newBeatsLogger(filename string, opts ...beatsLoggerOption) *beatsLogger {
	l := &beatsLogger{
		logFilename: filename,
		buffer:      &bytes.Buffer{},
	}
	for _, opt := range opts {
		opt(l)
	}
	log.Debugf("create beatsLogger to ( %s ", l.logFilename)
	return l
}

func (l *beatsLogger) Write(line []byte) (n int, err error) {
	log.Tracef("receveid line of ( %d ) bytes for beatsLogger", len(line))
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
	return len(line), nil
}

func (l *beatsLogger) Close() error {
	log.Tracef("closed beatsLogger")
	return nil
}
