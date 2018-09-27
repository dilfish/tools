package tools

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type AppendStruct struct {
	file *os.File
	c    chan os.Signal
	fn   string
	err  error
}

func openFile(fn string) (*os.File, error) {
	return os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
}

func (as *AppendStruct) wait() {
	signal.Notify(as.c, syscall.SIGUSR1)
	for {
		<-as.c
		io.WriteString(os.Stderr, "we got an signal, restart at "+TimeStr())
		f, err := openFile(as.fn)
		if err != nil {
			as.err = err
			log.Println("open file:", err)
			continue
		}
		as.file = f
	}
}

func NewAppender(fn string) (*AppendStruct, error) {
	var as AppendStruct
	f, err := openFile(fn)
	if err != nil {
		return nil, err
	}
	as.file = f
	as.fn = fn
	as.c = make(chan os.Signal)
	go as.wait()
	return &as, nil
}

func (as *AppendStruct) Close() {
	close(as.c)
}

func (as *AppendStruct) Write(bt []byte) (int, error) {
	if as.err != nil {
		return 0, as.err
	}
	n, err := as.file.Write(bt)
	if err != nil {
		as.err = err
	}
	return n, err
}

func Daemon() {
	os.Stdout.Close()
	os.Stdin.Close()
}

func InitLog(fn, prefix string) *log.Logger {
	as, err := NewAppender(fn)
	if err != nil {
		return nil
	}
	if prefix == "" {
		prefix = "DefAppendLogger "
	}
	if prefix[len(prefix)-1] != ' ' {
		prefix = prefix + " "
	}
	return log.New(as, prefix, log.LstdFlags|Lshortfile)
}
