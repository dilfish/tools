// sean at shanghai
// tcp proxy
// 2020

package main

import (
	"errors"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// ErrBadIP indicates ip is not valid
var ErrBadIP = errors.New("bad ip value")

// ErrBadPort indicates port is not valid or not authorized
var ErrBadPort = errors.New("bad port value")

// StatType in|out
type StatType string

const (
	// StatIn in bounds
	StatIn StatType = "in"
	// StatOut out bounds
	StatOut StatType = "out"
)

// Stat gets data every 2 minutes
type Stat struct {
	Inb  int64
	Outb int64
	Ts   time.Time
	Lock sync.Mutex
}

// TcpProxy stores ip and port information
type TcpProxy struct {
	LocalPort  int
	RemotePort int
	LocalIP    net.IP
	RemoteIP   net.IP
	Stat       *Stat
}

// NewProxy create a proxy of TCP
func NewProxy(localP, dstP int, localIP, dstIP string) (*TcpProxy, error) {
	lIP := net.ParseIP(localIP)
	dIP := net.ParseIP(dstIP)
	if lIP == nil || dIP == nil {
		return nil, ErrBadIP
	}
	if localP == 0 || dstP == 0 {
		return nil, ErrBadPort
	}
	var p TcpProxy
	p.Stat = &Stat{}
	p.LocalIP = lIP
	p.RemoteIP = dIP
	p.LocalPort = localP
	p.RemotePort = dstP
	return &p, nil
}

// GetNetClass judges we using tcp4 or tcp6
func GetNetClass(ip net.IP) string {
	if ip.To4() != nil {
		return "tcp4"
	}
	return "tcp6"
}

// Run runs a tcp proxy
func (p *TcpProxy) Run() error {
	var addr net.TCPAddr
	addr.Port = p.LocalPort
	ls, err := net.ListenTCP(GetNetClass(p.LocalIP), &addr)
	if err != nil {
		log.Println("listen tcp error:", err)
		return err
	}
	go p.StatCounter()
	for {
		c, err := ls.AcceptTCP()
		if err != nil {
			log.Println("accept error", err)
			continue
		}
		go p.Proxy(c)
	}
}

// Proxy starts a server and client
func (p *TcpProxy) Proxy(c *net.TCPConn) {
	defer c.Close()
	// log.Println("we get a conn from", c.RemoteAddr())
	var raddr net.TCPAddr
	raddr.IP = p.RemoteIP
	raddr.Port = p.RemotePort
	r, err := net.DialTCP(GetNetClass(p.RemoteIP), nil, &raddr)
	if err != nil {
		log.Println("dial remote error:", err)
		return
	}
	go p.LoopCopy(c, r, StatIn)
	p.LoopCopy(r, c, StatOut)
}

// LoopCopy copies data until error occur
func (p *TcpProxy) LoopCopy(dst io.Writer, src io.Reader, statType StatType) {
	en := 0
	for {
		n, err := io.Copy(dst, src)
		p.AddStat(n, statType)
		if err != nil {
			log.Println("io.Copy is", n, err)
		}
		if err != nil {
			en = en + 1
			if en > 2 {
				return
			}
		}
	}
}

// StatCounter reset stat every 2 minutes
func (p *TcpProxy) StatCounter() {
	c := time.Tick(time.Minute * 2)
	for range c {
		p.Stat.ClearOut()
	}
	return
}

// AddStat goes to stat.Add
func (p *TcpProxy) AddStat(n int64, statType StatType) {
	p.Stat.Add(n, statType)
}

// ClearOut clear all statistics and print them
func (stat *Stat) ClearOut() {
	var inb, outb int64
	stat.Lock.Lock()
	inb = stat.Inb
	outb = stat.Outb
	stat.Inb = 0
	stat.Outb = 0
	stat.Lock.Unlock()
	log.Println("current stat, in:", inb, " and out:", outb)
}

// Add add n to current stat
func (stat *Stat) Add(n int64, statType StatType) {
	stat.Lock.Lock()
	defer stat.Lock.Unlock()
	if statType == StatIn {
		stat.Inb = stat.Inb + n
	} else {
		stat.Outb = stat.Outb + n
	}
}
