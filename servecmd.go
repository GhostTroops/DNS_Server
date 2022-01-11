package main

import (
	"flag"
	"fmt"
	"log"
	"net"
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
		if 0 < index {
			return true
		}
	}
	return false
}

func parseQuery(m *dns.Msg, addressOfRequester net.Addr) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)
			if testIs(q.Name) {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}
}
func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	// func ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	// logrus.Infof("in")
	// msg := dns.Msg{}
	// dns.Responsewriter.RemoteAddr()
	addressOfRequester := w.RemoteAddr()
	msg := new(dns.Msg)
	msg.SetReply(r)
	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(msg, addressOfRequester)
	}

	// switch r.Question[0].Qtype {
	// case dns.TypeA:
	// 	msg.Authoritative = true
	// 	domain1 := msg.Question[0].Name
	// 	// logrus.Infof("domain: %v", domain1)
	// 	if testIs(domain1) {
	// 		logrus.Infof("domainxx:  %v %v  %v", addressOfRequester, domain1, ip)
	// 		msg.Answer = append(msg.Answer, &dns.A{
	// 			Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
	// 			A:   net.ParseIP(ip),
	// 		})
	// 	} else {
	// 		logrus.Infof("not do")
	// 	}

	// }
	w.WriteMsg(msg)
}

func main() {
	flag.StringVar(&domain, "domain", "51pwn.com", "set domain eg: 51pwn.com")
	flag.StringVar(&ip, "ip", "199.180.115.7", "set domain server ip, eg: 222.44.11.3")
	flag.StringVar(&resUrl, "resUrl", "http://127.0.0.1/dnsRecode", "Set the url that accepts dns parsing logs, eg: http://127.0.0.1/dnsRecode")
	flag.StringVar(&logLevel, "level", "INFO", "set loglevel, option")
	flag.Parse()
	a := regexp.MustCompile(`[,;]`)
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
		log.Fatalf("Failed to set udp listener %s\n", err.Error())
	}

}
