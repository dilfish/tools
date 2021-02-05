package http3

import (
	"errors"
	"fmt"
	"log"

	"github.com/miekg/dns"
)

var ErrInvalidName = errors.New("not a valid domain name")

func DoHTTPS(name, tp string) (*dns.MsgHdr, []dns.RR, []dns.RR, []dns.RR, error) {
	if name[len(name)-1] != '.' {
		name = name + "."
	}
	if len(name) < 3 {
		return nil, nil, nil, nil, ErrInvalidName
	}
	msg := new(dns.Msg)
	msg.Id = dns.Id()
	msg.RecursionDesired = true
	msg.Question = make([]dns.Question, 1)
	if tp == "https" {
		msg.Question[0] = dns.Question{
			Name:   name,
			Qtype:  dns.TypeHTTPS,
			Qclass: dns.ClassINET,
		}
	} else {
		msg.Question[0] = dns.Question{
			Name:   name,
			Qtype:  dns.TypeSVCB,
			Qclass: dns.ClassINET,
		}
	}
	c := new(dns.Client)
	in, _, err := c.Exchange(msg, "1.1.1.1:53")
	if err != nil {
		log.Println("Query 1.1.1.1 for", name, "error:", err)
		return nil, nil, nil, nil, err
	}
	return &in.MsgHdr, in.Answer, in.Ns, in.Extra, nil
}

func PrintResult(hdr *dns.MsgHdr, rr, ns, ex []dns.RR) {
	fmt.Println(hdr)
	if len(rr) > 0 {
		fmt.Println("DNS Answers:")
	}
	for _, r := range rr {
		h, ok := r.(*dns.HTTPS)
		if ok {
			fmt.Println(h.Hdr.Name, h.Hdr.Ttl, dns.Type(h.Hdr.Rrtype).String())
			for _, v := range h.Value {
				fmt.Println(v.Key().String() + ": " + v.String())
			}
		} else {
			s, ok := r.(*dns.SVCB)
			if ok {
				fmt.Println(s.Hdr.Name, s.Hdr.Ttl, dns.Type(s.Hdr.Rrtype).String())
				fmt.Println("pri and target:", s.Priority, s.Target)
				for _, v := range s.Value {
					fmt.Println(v.Key(), v.String())
				}
			} else {
				fmt.Println("Not https nor scvb:", r)
			}
		}
	}
	if len(ns) > 0 {
		fmt.Println("DNS Authoratives:")
		for _, n := range ns {
			fmt.Println(n)
		}
	}
	if len(ex) > 0 {
		fmt.Println("Extra Data:")
		for _, e := range ex {
			fmt.Println(e)
		}
	}
}
