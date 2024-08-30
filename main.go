package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"fmt"
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
	server          = "223.5.5.5"
	accountID       string
	accessKeySecret string
	accessKeyID     string
	listenAddress   = ":53"
	timeout         time.Duration = 3 * time.Second
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "ArashiDNS-Aha",
		Short: "阿里云递归（公共）HTTP DNS 客户端",
		Run: func(cmd *cobra.Command, args []string) {
			dns.HandleFunc(".", handleDNSQuery)
			var servers []*dns.Server
			servers = append(servers, UDPDNS(listenAddress))
			servers = append(servers, TCPDNS(listenAddress))
			fmt.Printf("Now listening on: %s\n", listenAddress)
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
		},
	}

	rootCmd.Flags().StringVar(&accountID, "accountID", "", "云解析-公共 DNS 控制台的 Account ID")
	rootCmd.MarkFlagRequired("accountID")
	rootCmd.Flags().StringVar(&accessKeySecret, "accessKeySecret", "", "云解析-公共 DNS 控制台创建密钥中的 AccessKey 的 Secret")
	rootCmd.MarkFlagRequired("accessKeySecret")
	rootCmd.Flags().StringVar(&accessKeyID, "accessKeyID", "", "云解析-公共 DNS 控制台创建密钥中的 AccessKey 的 ID")
	rootCmd.MarkFlagRequired("accessKeyID")
	rootCmd.Flags().StringVar(&server, "server", "223.5.5.5", "设置的服务器的地址")
	rootCmd.Flags().StringVar(&listenAddress, "listen", ":53", "监听的地址与端口")
	rootCmd.Flags().DurationVar(&timeout, "timeout", 3*time.Second, "等待回复的超时时间")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func handleDNSQuery(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	msg.Authoritative = true

	for _, q := range r.Question {
		switch q.Qtype {
		case dns.TypeA, dns.TypeAAAA, dns.TypeCNAME, dns.TypeNS, dns.TypeTXT:
			answer, err := queryHTTPDNS(q.Name, dns.Type(q.Qtype).String())
			if err != nil || answer.Status != 0 {
				msg.Rcode = dns.RcodeServerFailure
			} else {
				fmt.Println(answer)
				for _, ans := range answer.Answer {
					record := getDNSRecord(ans)
					msg.Answer = append(msg.Answer, record)
				}
			}
		}
	}

	w.WriteMsg(&msg)
}

func queryHTTPDNS(name, qtype string) (*DNSEntity, error) {
	ts := fmt.Sprintf("%d", time.Now().Unix())
	key := sha256.Sum256([]byte(accountID + accessKeySecret + ts + name + accessKeyID))
	keyStr := hex.EncodeToString(key[:])
	url := fmt.Sprintf("http://%s/resolve?name=%s&type=%s&uid=%s&ak=%s&key=%s&ts=%s", server, name, qtype, accountID, accessKeyID, keyStr, ts)
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
		rr.Txt = []string{ans.Data}
		return rr
	}
}
