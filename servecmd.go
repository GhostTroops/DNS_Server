package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
)

var (
	domain, ip, resUrl, logLevel string
	aDomain                      []string
)

type handler struct{}

func testIs(s1 string) bool {
	for _, s := range aDomain {
		index := strings.Index(s1, "."+s)
		if 0 <= index {
			return true
		}
	}
	return false
}

// 解决相同多次请求的问题
func sendReq(addressOfRequester net.Addr, domain1 string) {
	ip1 := fmt.Sprintf("%s", addressOfRequester)
	a := regexp.MustCompile(`[:]`)
	ip1 = a.Split(ip1, -1)[0]
	i := strings.Count(domain1, "") - 2
	domain1 = domain1[0:i]
	logrus.Info(domain1 + " " + ip1)
	post_body := bytes.NewReader([]byte(fmt.Sprintf(`{"ip":"%s","domain":"%s"}`, ip1, domain1)))
	go http.Post(resUrl, "application/json", post_body)
}

func parseQuery(m *dns.Msg, addressOfRequester net.Addr) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:

			if testIs(q.Name) {
				// logrus.Info("Query for %s %v\n", q.Name, addressOfRequester)
				go sendReq(addressOfRequester, q.Name)
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}
}

func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	addressOfRequester := w.RemoteAddr()
	msg := new(dns.Msg)
	msg.SetReply(r)
	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(msg, addressOfRequester)
	}

	w.WriteMsg(msg)
}

func main() {
	flag.StringVar(&domain, "domain", "51pwn.com", "set domain eg: 51pwn.com")
	flag.StringVar(&ip, "ip", "199.180.115.7", "set domain server ip, eg: 222.44.11.3")
	flag.StringVar(&resUrl, "resUrl", "http://127.0.0.1/dnsRecode", "Set the url that accepts dns parsing logs, eg: http://127.0.0.1/dnsRecode")
	flag.StringVar(&logLevel, "level", "INFO", "set loglevel, option")
	flag.Parse()
	a := regexp.MustCompile(`[,;:]`)
	aDomain = a.Split(domain, -1)

	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	default:
		logrus.SetLevel(logrus.WarnLevel)
	}
	// dns.HandleFunc(domain+".", ServeDNS)
	srv := &dns.Server{Addr: ":" + strconv.Itoa(53), Net: "udp"}

	srv.Handler = &handler{}
	// srv.Handler = &ServeDNS{}
	err := srv.ListenAndServe()
	defer srv.Shutdown()
	if err != nil {
		logrus.Fatalf("Failed to set udp listener %s\n", err.Error())
	}

}
