package main

import (
	"bytes"
	"encoding/json"
)

type _Options struct {
	API    APIOptions `json:"api"`
	Server string     `json:"server"`
	DNS    DNSOptions `json:"dns"`
	TLS    TLSOptions `json:"tls"`
}

type Options _Options

func (o *Options) UnmarshalJSON(content []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	err := decoder.Decode((*_Options)(o))
	if err != nil {
		return err
	}
	if o.Server == "" {
		o.Server = "223.5.5.5"
	}
	return nil
}

type APIOptions struct {
	AccountID       string       `json:"account_id"`
	AccessKeyID     string       `json:"access_key_id"`
	AccessKeySecret string       `json:"access_key_secret"`
	ExtraOptions    ExtraOptions `json:"extra"`
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
