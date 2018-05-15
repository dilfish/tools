package tools

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
	"errors"
)


var ErrBadFmt = errors.New("bad format")
var ErrNoSuch = errors.New("no such")
var ErrDupData = errors.New("dup data")


func RandInt(w int) int32 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int31n(int32(w))
}

func RandStr(w int) string {
	rand.Seed(time.Now().UnixNano())
	base := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	str := ""
	for i := 0; i < w; i++ {
		idx := rand.Int31n(int32(len(base)))
		str = str + string(base[idx])
	}
	return str
}

func ReadFile(fn string) ([]byte, error) {
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

func DoPost(url string, v *url.Values) ([]byte, error) {
	resp, err := http.PostForm(url, *v)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func DoGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

type LineFunc func(line string) error

func ReadLine(fn string, lf LineFunc) error {
	file, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer file.Close()
	rd := bufio.NewReader(file)
	for {
		line, err := rd.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}
		if line == "" {
			continue
		}
		line = line[:len(line)-1]
		err = lf(line)
		if err != nil {
			return err
		}
	}
	return nil
}

func Proxy(c *net.TCPConn) {
	defer c.Close()
	now := time.Now()
	fmt.Println(now, "we get an conn from", c.RemoteAddr())
	fmt.Println(now, "and we are going to 119.28.77.61:8000...")
	var raddr net.TCPAddr
	raddr.IP = net.ParseIP("119.28.77.61")
	raddr.Port = 8000
	r, err := net.DialTCP("tcp4", nil, &raddr)
	if err != nil {
		fmt.Println("dial remote", err)
		return
	}
	go io.Copy(c, r)
	io.Copy(r, c)
}

func Run() error {
	var addr net.TCPAddr
	addr.Port = 8080
	ls, err := net.ListenTCP("tcp4", &addr)
	if err != nil {
		return err
	}
	for {
		c, err := ls.AcceptTCP()
		if err != nil {
			fmt.Println("accept error", err)
			continue
		}
		go Proxy(c)
	}
}
