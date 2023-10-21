/*+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

This package implements the shell's wrapper.

+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/
package business

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	// creack (>= 1.1.10) solves the error "panic: fork/exec /bin/bash: Setctty set but Ctty not valid in child"
	"github.com/creack/pty"
	log "github.com/sirupsen/logrus"
	"golang.org/x/term"
)

const (
	bashShellCmd = "/bin/bash"
	dumpLoggerDefaultFile = "/var/log/s-h-entinel.log"
	jsonLoggerDefaultFile = "/var/log/s-h-entinel.json"
)

type logger interface {
	Write(p []byte) (n int, err error)
	Close() error
}

type Sentinel struct {
	buffer      *bytes.Buffer
	loggers     []logger
	shell       string
}

type SentinelOption func(*Sentinel)

func WithJsonLogger(enabled bool, filename string) SentinelOption {
	return func(s *Sentinel) {
		if enabled {
			s.loggers = append(s.loggers, newJsonLogger(filename))
		}
	}
}

func WithDumpLogger(enabled bool, filename string) SentinelOption {
	return func(s *Sentinel) {
		if enabled {
			s.loggers = append(s.loggers, newDumpLogger(filename))
		}
	}
}

func New(opts ...SentinelOption) *Sentinel {
	s := &Sentinel{
		buffer:      &bytes.Buffer{},
		shell: bashShellCmd,
		loggers: []logger{},
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Write: this func splits the given array of bytes into line of text, 
//        then it writes every line of text to every defined logger.
func (w *Sentinel) Write(p []byte) (n int, err error) {
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

// Start: it starts a shell, getting the stdin from that, 
//        then it write every byte from the stdin to a multiwriter,
//        the latter includes the os.Stdout (to make the shell functioning for the user) 
//        and the Sentinel itself, for auditing purpose.
func (w *Sentinel) Start() {
	signal.Notify(channel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	go signalsListener()
	cmd := exec.Command(w.shell)
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
	exitcode := <-exitchannel
	os.Exit(exitcode)
}
