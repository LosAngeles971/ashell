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

type Wrapper struct {
	l *logger
}

func New() *Wrapper {
	return &Wrapper{
		l: newLogger(withLogFile("/tmp/ashell.log"), withAuditFile("/tmp/ashell.json")),
	}
}

func (w *Wrapper) Start() {
	defer w.l.Close()
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
	mw := io.MultiWriter(os.Stdout, w.l)
	go io.Copy(tty, os.Stdin)
	io.Copy(mw, tty)
}
