package business

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log/syslog"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	sysLogTag      = "ashell"
	sysLogProtocol = "udp"
)

// A log entry can refer to the execution of a test suite or a test case
type CommandEvent struct {
	Timestamp string `json:"@timestamp"`
	Command   string `json:"command"`
	Output    string `json:"output,omitempty"`
}

type logger struct {
	logFileEnabled   bool
	auditFileEnabled bool
	sysLogEnabled    bool
	logFilename      string
	auditFilename    string
	sysLogTarget     string
	logFile          *os.File
	sysLog           *syslog.Writer
	buffer           *bytes.Buffer
}

type loggerOption func(*logger)

func withLogFile(filename string) loggerOption {
	return func(l *logger) {
		var err error
		l.logFilename = filename
		l.logFile, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0775)
		if err != nil {
			panic(err)
		}
		l.logFileEnabled = true
	}
}

func withAuditFile(filename string) loggerOption {
	return func(l *logger) {
		l.auditFilename = filename
		l.auditFileEnabled = true
	}
}

func withSysLog(syslogTarget string) loggerOption {
	return func(l *logger) {
		var err error
		l.sysLogTarget = syslogTarget
		l.sysLog, err = syslog.Dial(sysLogProtocol, syslogTarget, syslog.LOG_DEBUG, sysLogTag)
		if err != nil {
			panic(err)
		}
		l.sysLogEnabled = true
	}
}

func newLogger(opts ...loggerOption) *logger {
	l := &logger{
		logFileEnabled: false,
		sysLogEnabled:  false,
		buffer:         &bytes.Buffer{},
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// logEvent logs a command event to a file
func (l *logger) logEvent(line []byte) {
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
	f, err := os.OpenFile(l.auditFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0775)
	if err != nil {
		log.Error(err)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(string(j) + "\n"); err != nil {
		log.Error(err)
	}
}

func (l *logger) write(p []byte) (n int, err error) {
	n = len(p)
	l.buffer.Write(p)
	if bytes.Contains(p, []byte{'\n'}) {
		for _, line := range bytes.Split(l.buffer.Bytes(), []byte{'\n'}) {
			if len(line) == 0 {
				continue
			}
			if l.logFileEnabled {
				l.logFile.Write(line)
			}
			if l.auditFileEnabled {
				l.logEvent(line)
			}
			if l.sysLogEnabled {
				_, err := l.sysLog.Write(line)
				if err != nil {
					log.Errorf("error writing to syslog ( %s ) - %v", l.sysLogTarget, err)
				}
			}
		}
		l.buffer.Reset()
	}
	return n, nil
}

func (l *logger) Write(p []byte) (n int, err error) {
	return l.write(p)
}

func (l *logger) Close() {
	if l.logFileEnabled {
		l.logFile.Close()
	}
	if l.sysLogEnabled {
		l.sysLog.Close()
	}
}
