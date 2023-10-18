package business

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type dumpLogger struct {
	logFilename string
}

type dumpLoggerOption func(*dumpLogger)

func newDumpLogger(filename string, opts ...dumpLoggerOption) *dumpLogger {
	l := &dumpLogger{
		logFilename: filename,
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l *dumpLogger) Write(line []byte) (n int, err error) {
	log.Tracef("receveid line of ( %d ) bytes for dumpLogger", len(line))
	if f, err := os.OpenFile(l.logFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0775); err != nil {
		log.Errorf("error auditing commands to ( %s ) - %v", l.logFilename, err)
	} else {
		f.WriteString(fmt.Sprintf("%s\n", string(line)))
		f.Close()
	}
	return len(line), nil
}

func (l *dumpLogger) Close() error {
	log.Tracef("closed dumpLogger")
	return nil
}
