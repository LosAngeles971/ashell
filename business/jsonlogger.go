/*+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

This package implements the JSON logger.

+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/

package business

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// LogEvent: this struct represent a single entry into the JSON log file.
type LogEvent struct {
	Timestamp string `json:"@timestamp"`
	Command   string `json:"command"`
	Output    string `json:"output,omitempty"`
}

// jsonLogger: it implements the JSON logger
type jsonLogger struct {
	logFilename string // JSON target file
}

// beatsLoggerOption: it defines the function to optionally use when you create a JSON logger
type beatsLoggerOption func(*jsonLogger)

// newJsonLogger: it returns a new instance of JSON logger
func newJsonLogger(filename string, opts ...beatsLoggerOption) *jsonLogger {
	l := &jsonLogger{
		logFilename: filename,
	}
	for _, opt := range opts {
		opt(l)
	}
	log.Debugf("create beatsLogger to ( %s ", l.logFilename)
	return l
}

// Write: it appends a new line as LogEvent into the target JSON log file
func (l *jsonLogger) Write(line []byte) (n int, err error) {
	log.Tracef("receveid line of ( %d ) bytes for beatsLogger", len(line))
	event := LogEvent{
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

// Close: it is a useless function implemented only to be consistent with the Writer interface
func (l *jsonLogger) Close() error {
	log.Tracef("closed beatsLogger")
	return nil
}
