package business

import (
	"bytes"
	"log/syslog"
	log "github.com/sirupsen/logrus"
)

const (
	sysLogTag      = "ashell"
	sysLogProtocol = "udp"
)

type syslogLogger struct {
	sysLogTarget string
	sysLog       *syslog.Writer
	buffer       *bytes.Buffer
}

type osLoggerOption func(*syslogLogger)

func newSyslogLogger(syslogTarget string, opts ...osLoggerOption) *syslogLogger {
	var err error
	l := &syslogLogger{
		sysLogTarget: syslogTarget,
	}
	l.sysLog, err = syslog.Dial(sysLogProtocol, syslogTarget, syslog.LOG_DEBUG, sysLogTag)
	if err != nil {
		panic(err)
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l *syslogLogger) Write(p []byte) (n int, err error) {
	log.Tracef("receveid line of ( %d ) bytes for dumpLogger", len(line))
	if _, err := l.sysLog.Write(line); err != nil {
		log.Errorf("error auditing commands to ( %s ) - %v", l.logFilename, err)
		return len(line), err
	}
	return len(line), nil
}

func (l *syslogLogger) Close() error {
	return l.sysLog.Close()
}
