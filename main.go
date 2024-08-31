package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/miekg/dns"
	"github.com/spf13/cobra"
)

var (
	configPath      string
	workerDir       string
	options         = &Options{}
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "AhaDNS",
		Short: "阿里云递归（公共）HTTP DNS 客户端",
		Run: func(cmd *cobra.Command, args []string) {
			err := run()
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "配置文件")
	rootCmd.MarkFlagRequired("config")
	rootCmd.Flags().StringVarP(&workerDir, "directory", "D", "", "工作目录")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	options, err := readConfig(configPath)
	if err != nil {
		return err
	}
	if workerDir != "" {
		_, err = os.Stat(workerDir)
		if err != nil {
			os.Mkdir(workerDir, 0o777)
		}
		err = os.Chdir(workerDir)
		if err != nil {
			return err
		}
	}

	dns.HandleFunc(".", handleDNSQuery)
	var servers []*dns.Server
	if options.DNS.UDPOption.Enabled {
		if server, err := UDPDNS(options.DNS.UDPOption.Listen, options.DNS.UDPOption.ListenPort); err != nil {
			return fmt.Errorf("Failed to start UDP server: %v\n", err)
		} else {
			servers = append(servers, server)
		}
	}
	if options.DNS.TCPOption.Enabled {
		if server, err := TCPDNS(options.DNS.TCPOption.Listen, options.DNS.TCPOption.ListenPort); err != nil {
			return fmt.Errorf("Failed to start TCP server: %v\n", err)
		} else {
			servers = append(servers, server)
		}
	}
	if options.DNS.TLSOption.Enabled {
		if !options.TLS.Enabled {
			return fmt.Errorf("TLS options are not enabled")
		}
		if options.TLS.CertPath == "" {
			return fmt.Errorf("TLS certificate path is not set")
		}
		if options.TLS.KeyPath == "" {
			return fmt.Errorf("TLS key path is not set")
		}
		if server, err := TLSDNS(options.DNS.TLSOption.Listen, options.DNS.TLSOption.ListenPort, options.TLS.CertPath, options.TLS.KeyPath); err != nil {
			return fmt.Errorf("Failed to start TLS server: %v\n", err)
		} else {
			servers = append(servers, server)
		}
	}

	fmt.Println("Application started. Press Ctrl+C to shut down.")

	for _, server := range servers {
		go func() {
			if err := server.ListenAndServe(); err != nil {
				log.Fatalf("Failed to start server: %v\n", err)
			}
		}()
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	for _, server := range servers {
		server.Shutdown()
	}
	return nil
}

func readConfig(path string) (*Options, error) {
	var (
		configContent []byte
		err           error
	)
	if path == "stdin" {
		configContent, err = io.ReadAll(os.Stdin)
	} else {
		configContent, err = os.ReadFile(path)
	}
	if err != nil {
		return nil, fmt.Errorf("%v: read config at %s", err, path)
	}
	err = options.UnmarshalJSON(configContent)
	if err != nil {
		return nil, fmt.Errorf("%v: decode config at %s", err, path)
	}
	return options, nil
}

func JoinIPPort[T int | uint16](ip string, port T) (string, error) {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return "", fmt.Errorf("invalid IP address")
	}

	if parsedIP.To4() == nil { // 是 IPv6 地址
		return fmt.Sprintf("[%s]:%d", ip, port), nil
	}

	return fmt.Sprintf("%s:%d", ip, port), nil
}

func handleDNSQuery(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	msg.Authoritative = true

	for _, q := range r.Question {
		switch q.Qtype {
		case dns.TypeA, dns.TypeAAAA, dns.TypeCNAME, dns.TypeNS, dns.TypeTXT:
			answer, err := queryHTTPDNS(options, q.Name, dns.Type(q.Qtype).String())
			if err != nil || answer.Status != 0 {
				msg.Rcode = dns.RcodeServerFailure
			} else {
				for _, ans := range answer.Answer {
					record := getDNSRecord(ans)
					msg.Answer = append(msg.Answer, record)
				}
			}
		case dns.TypeMX:
			answer, err := queryRawDNS(options, q.Name)
			if err != nil {
				msg.Rcode = dns.RcodeServerFailure
			} else {
				for _, ans := range answer.Answer {
					msg.Answer = append(msg.Answer, ans)
				}
			}
		}
	}

	w.WriteMsg(&msg)
}

func queryRawDNS(options *Options, name string) (*dns.Msg, error) {
	dnsClient := new(dns.Client)

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(name), dns.TypeMX)
	msg.RecursionDesired = true

	serverAddress, err := JoinIPPort(options.Server.Address, options.Server.UDPPort)
	if err != nil {
		return nil, err
	}
	r, _, err := dnsClient.Exchange(msg, serverAddress)
	if err != nil {
		return nil, fmt.Errorf("Failed to exchange: %v", err)
	}

	if r.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("Query failed: %s", dns.RcodeToString[r.Rcode])
	}
	return r, nil
}

func queryHTTPDNS(options *Options, name string, qtype string) (*DNSEntity, error) {
	ts := fmt.Sprintf("%d", time.Now().Unix())
	key := sha256.Sum256([]byte(options.API.AccountID + options.API.AccessKeySecret + ts + name + options.API.AccessKeyID))
	keyStr := hex.EncodeToString(key[:])
	url := fmt.Sprintf("http://%s/resolve?name=%s&type=%s&uid=%s&ak=%s&key=%s&ts=%s", options.Server.Address, name, qtype, options.API.AccountID, options.API.AccessKeyID, keyStr, ts)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result DNSEntity
	err = json.Unmarshal(body, &result)
	return &result, err
}

func getDNSRecord(ans Answer) dns.RR {
	header := dns.RR_Header{Name: ans.Name, Rrtype: uint16(ans.Type), Class: dns.ClassINET, Ttl: uint32(ans.TTL)}
	switch ans.Type {
	case 1: // A
		rr := new(dns.A)
		rr.Hdr = header
		rr.A = net.ParseIP(ans.Data)
		return rr
	case 28: // AAAA
		rr := new(dns.AAAA)
		rr.Hdr = header
		rr.AAAA = net.ParseIP(ans.Data)
		return rr
	case 5: // CNAME
		rr := new(dns.CNAME)
		rr.Hdr = header
		rr.Target = ans.Data
		return rr
	case 2: // NS
		rr := new(dns.NS)
		rr.Hdr = header
		rr.Ns = ans.Data
		return rr
	case 16: // TXT
		rr := new(dns.TXT)
		rr.Hdr = header
		cleanedData := strings.Trim(ans.Data, "\"") // 去掉引号
		rr.Txt = []string{cleanedData}
		return rr
	default:
		rr := new(dns.TXT)
		rr.Hdr = header
		cleanedData := strings.Trim(ans.Data, "\"") // 去掉引号
		rr.Txt = []string{cleanedData}
		return rr
	}
}
