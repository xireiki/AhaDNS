package main

import (
	"crypto/tls"
	"github.com/miekg/dns"
	"strconv"
)

func UDPDNS(Listen string, ListenPort int) (*dns.Server, error) {
	return &dns.Server{
		Addr: Listen + ":" + strconv.Itoa(ListenPort),
		Net: "udp",
	}, nil
}

func TCPDNS(Listen string, ListenPort int) (*dns.Server, error) {
	return &dns.Server{
		Addr: Listen + ":" + strconv.Itoa(ListenPort),
		Net: "tcp",
	}, nil
}

func TLSDNS(Listen string, ListenPort int, certPath string, keyPath string) (*dns.Server, error) {
	dnsServer := &dns.Server{
		Addr:    Listen + ":" + strconv.Itoa(ListenPort),
		Net:     "tcp-tls",
	}
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS13,
		Certificates: []tls.Certificate{cert},
	}
	dnsServer.TLSConfig = tlsConfig
	return dnsServer, nil
}
