package tools

import (
	"net"
)

func IP2Num(ipstr string) uint32 {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	if ip == nil {
		return 0
	}
	n := uint32(0)
	n = n + uint32(ip[0])*256*256*256
	n = n + uint32(ip[1])*256*256
	n = n + uint32(ip[2])*256
	n = n + uint32(ip[3])
	return n
}

func Num2IP(ipnum uint32) string {
	c1 := ipnum / 256 / 256 / 256
	c2 := (ipnum / 256 / 256) % 256
	c3 := (ipnum / 256) % 256
	c4 := ipnum % 256
	return net.IPv4(byte(c1), byte(c2), byte(c3), byte(c4)).String()
}
