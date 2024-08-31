package main

import (
	"bytes"
	"encoding/json"
)

type _Options struct {
	API    APIOptions    `json:"api"`
	Server ServerOptions `json:"server"`
	DNS    DNSOptions    `json:"dns"`
	TLS    TLSOptions    `json:"tls"`
}

type Options _Options

func (o *Options) UnmarshalJSON(content []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	err := decoder.Decode((*_Options)(o))
	if err != nil {
		return err
	}
	o.ServerOptionDefaultCheck()
	return nil
}

func (o *Options) ServerOptionDefaultCheck() error {
	if o.Server.Address == "" {
		o.Server.Address = "223.5.5.5"
	}
	if o.Server.UDPPort == 0 {
		o.Server.UDPPort = 53
	}
	if o.Server.TCPPort == 0 {
		o.Server.TCPPort = 53
	}
	if o.Server.TLSPort == 0 {
		o.Server.TLSPort = 853
	}
	if o.Server.HTTPPort == 0 {
		o.Server.HTTPPort = 80
	}
	if o.Server.HTTPSPort == 0 {
		o.Server.HTTPSPort = 443
	}
	return nil
}

type APIOptions struct {
	AccountID       string       `json:"account_id"`
	AccessKeyID     string       `json:"access_key_id"`
	AccessKeySecret string       `json:"access_key_secret"`
	ExtraOptions    ExtraOptions `json:"extra"`
}

type ServerOptions struct {
	Address   string `json:"address"`
	UDPPort   uint16 `json:"udp_port"`
	TCPPort   uint16 `json:"tcp_port"`
	TLSPort   uint16 `json:"tls_port"`
	HTTPPort  uint16 `json:"http_port"`
	HTTPSPort uint16 `json:"https_port"`
}

type ExtraOptions struct {
	Enabled       bool   `json:"enabled"`
	RequestMethod string `json:"method"` // udp/tcp/tls/http/https
}

type DNSOptions struct {
	UDPOption   ListenOptions `json:"udp"`
	TCPOption   ListenOptions `json:"tcp"`
	TLSOption   ListenOptions `json:"tls"`
	HTTPOption  ListenOptions `json:"http"`
	HTTPSOption ListenOptions `json:"https"`
}

type ListenOptions struct {
	Enabled    bool   `json:"enabled"`
	Listen     string `json:"listen"`
	ListenPort uint16 `json:"listen_port"`
}

type TLSOptions struct {
	Enabled  bool   `json:"enabled"`
	CertPath string `json:"cert_path"`
	KeyPath  string `json:"key_path"`
}
