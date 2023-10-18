package business

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	// creack (>= 1.1.10) solves the error "panic: fork/exec /bin/bash: Setctty set but Ctty not valid in child"
	"github.com/creack/pty"
	log "github.com/sirupsen/logrus"
	"golang.org/x/term"
)

const (
	shellCmd = "/bin/bash"
)

type logger interface {
	Write(p []byte) (n int, err error)
	Close() error
}

type Wrapper struct {
	buffer_size int
	buffer      *bytes.Buffer
	loggers     []logger
}

func New() *Wrapper {
	return &Wrapper{
		buffer_size: 1000,
		buffer:      &bytes.Buffer{},
		loggers: []logger{
			newDumpLogger("/tmp/ashell.log"),
			newBeatsLogger("/tmp/ashell.json"),
		},
	}
}

func (w *Wrapper) Write(p []byte) (n int, err error) {
	n = len(p)
	log.Tracef("receveid ( %d ) bytes for multi-writer", n)
	w.buffer.Write(p)
	if bytes.Contains(p, []byte{'\n'}) {
		for _, line := range bytes.Split(w.buffer.Bytes(), []byte{'\n'}) {
			if len(line) == 0 {
				continue
			}
			for _, l := range w.loggers {
				l.Write(line)
			}
		}
		w.buffer.Reset()
	}
	return n, nil
}

func (w *Wrapper) Start() {
	cmd := exec.Command(shellCmd)
	tty, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}
	defer tty.Close()
	previousState, err := term.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer term.Restore(0, previousState)
	go io.Copy(tty, os.Stdin)
	mw := io.MultiWriter(os.Stdout, w)
	io.Copy(mw, tty)
	log.Debugf("closing loggers...")
	for _, l := range w.loggers {
		l.Close()
	}
	log.Debugf("closed loggers")
}
