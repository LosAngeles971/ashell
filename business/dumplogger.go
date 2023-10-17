package business

import (
	"bytes"
	"os"
)

type dumpLogger struct {
	logFilename      string
	logFile          *os.File
	buffer           *bytes.Buffer
}

type dumpLoggerOption func(*dumpLogger)

func newDumpLogger(filename string, opts ...dumpLoggerOption) *dumpLogger {
	var err error
	l := &dumpLogger{
		logFilename: filename,
		buffer:         &bytes.Buffer{},
	}
	l.logFile, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0775)
	if err != nil {
		panic(err)
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}


func (l *dumpLogger) Write(p []byte) (n int, err error) {
	n = len(p)
	l.buffer.Write(p)
	if bytes.Contains(p, []byte{'\n'}) {
		for _, line := range bytes.Split(l.buffer.Bytes(), []byte{'\n'}) {
			if len(line) == 0 {
				continue
			}
			l.logFile.Write(line)
		}
		l.buffer.Reset()
	}
	return n, nil
}

func (l *dumpLogger) Close() error {
	return l.logFile.Close()
}
