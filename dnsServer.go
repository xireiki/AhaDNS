package main

import (
	"github.com/miekg/dns"
)

func UDPDNS(ListenAddress string) *dns.Server {
	return &dns.Server{Addr: ListenAddress, Net: "udp"}
}

func TCPDNS(ListenAddress string) *dns.Server {
	return &dns.Server{Addr: ListenAddress, Net: "tcp"}
}
