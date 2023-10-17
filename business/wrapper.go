package business

import (
	"io"
	"os"
	"os/exec"

	// creack (>= 1.1.10) solves the error "panic: fork/exec /bin/bash: Setctty set but Ctty not valid in child"
	"github.com/creack/pty"
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
	loggers []logger
}

func New() *Wrapper {
	return &Wrapper{
		buffer_size: 1000,
		loggers: []logger{
			newDumpLogger("/tmp/ashell.log"),
			newBeatsLogger("/tmp/ashell.json"),
		},
	}
}

func (w *Wrapper) multicopy(source *os.File) {
    var buffer = make([]byte, w.buffer_size)
    for {
		n, err := source.Read(buffer)
		if err != nil {
			panic(err)
		}
		if n == 0 {
			return
		}
		p := []byte{}
		for i := 0; i < n; i++ {
			p[i] = buffer[i]
		}
		for _, l := range w.loggers {
			l.Write(p)
		}
	}
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
	mw := io.MultiWriter(os.Stdout, w.loggers[0])
	go io.Copy(tty, os.Stdin)
	io.Copy(mw, tty)
	for _, l := range w.loggers {
		l.Close()
	}
}
