package business

import (
	"bytes"
	"log/syslog"
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
	n = len(p)
	l.buffer.Write(p)
	if bytes.Contains(p, []byte{'\n'}) {
		for _, line := range bytes.Split(l.buffer.Bytes(), []byte{'\n'}) {
			if len(line) == 0 {
				continue
			}
			_, err := l.sysLog.Write(line)
			if err != nil {
				panic(err)
			}
		}
		l.buffer.Reset()
	}
	return n, nil
}

func (l *syslogLogger) Close() error {
	return l.sysLog.Close()
}
