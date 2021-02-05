// Copyright 2018 Sean.ZH

package logh

import (
	"log"
	"os"
	"sync"
)

// AppendStruct holds log file needed
type AppendStruct struct {
	file   *os.File
	cClose chan bool
	fn     string
	err    error
	pid    int
	lock   sync.Mutex
}

func openFile(fn string) (*os.File, error) {
	return os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
}

func (as *AppendStruct) Restart() error {
	f, err := openFile(as.fn)
	if err != nil {
		as.err = err
		log.Println("open file:", as.fn, err)
		return err
	}
	as.lock.Lock()
	if as.file != nil {
		as.file.Close()
	}
	as.file = f
	as.lock.Unlock()
	return nil
}

// NewAppender create an append only log file with debug info
func NewAppender(fn string) (*AppendStruct, error) {
	var as AppendStruct
	as.fn = fn
	err := as.Restart()
	if err != nil {
		return nil, err
	}
	as.pid = os.Getpid()
	return &as, nil
}

// Close release all resources it holds
func (as *AppendStruct) Close() {
	as.lock.Lock()
	defer as.lock.Unlock()
	as.file.Close()
	as.file = nil
}

// Write api for file write
func (as *AppendStruct) Write(bt []byte) (int, error) {
	if as.err != nil {
		return 0, as.err
	}
	as.lock.Lock()
	defer as.lock.Unlock()
	n, err := as.file.Write(bt)
	if err != nil {
		as.err = err
	}
	return n, err
}

// Daemon close stdin stdout
func Daemon() {
	os.Stdout.Close()
	os.Stdin.Close()
	os.Stdout = nil
	os.Stdin = nil
}

// InitLog create a new log object
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
	return log.New(as, prefix, log.LstdFlags|log.Lshortfile)
}

// SetLog normal log set
func SetLog() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
