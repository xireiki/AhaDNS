package main

import (
	"crypto/tls"
	"github.com/miekg/dns"
)

func UDPDNS(Listen string, ListenPort uint16) (*dns.Server, error) {
	address, err := JoinIPPort(Listen, int(ListenPort))
	if err != nil {
		return nil, err
	}
	return &dns.Server{
		Addr: address,
		Net:  "udp",
	}, nil
}

func TCPDNS(Listen string, ListenPort uint16) (*dns.Server, error) {
	address, err := JoinIPPort(Listen, int(ListenPort))
	if err != nil {
		return nil, err
	}
	return &dns.Server{
		Addr: address,
		Net: "tcp",
	}, nil
}

func TLSDNS(Listen string, ListenPort uint16, certPath string, keyPath string) (*dns.Server, error) {
	address, err := JoinIPPort(Listen, int(ListenPort))
	if err != nil {
		return nil, err
	}
	dnsServer := &dns.Server{
		Addr: address,
		Net:  "tcp-tls",
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
