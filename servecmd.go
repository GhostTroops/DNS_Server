package main

import (
	"flag"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
)

var (
	domain, ip, resUrl, logLevel string
)

type handler struct{}

// func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
func ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	logrus.Infof("in")
	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain1 := msg.Question[0].Name
		// logrus.Infof("domain: %v", domain1)
		index := strings.Index(domain1, "."+domain)
		if 0 < index {
			// logrus.Infof("domainxx: %v", ip)
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(ip),
			})
		}
	}
	w.WriteMsg(&msg)
}

func main() {
	flag.StringVar(&domain, "domain", "51pwn.com", "set domain eg: 51pwn.com")
	flag.StringVar(&ip, "ip", "199.180.115.7", "set domain server ip, eg: 222.44.11.3")
	flag.StringVar(&resUrl, "resUrl", "http://127.0.0.1/dnsRecode", "Set the url that accepts dns parsing logs, eg: http://127.0.0.1/dnsRecode")
	flag.StringVar(&logLevel, "level", "INFO", "set loglevel, option")
	flag.Parse()

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
	dns.HandleFunc(domain+".", ServeDNS)
	srv := &dns.Server{Addr: ":" + strconv.Itoa(53), Net: "udp"}

	// srv.Handler = &handler{}
	// srv.Handler = &ServeDNS{}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %s\n", err.Error())
	}
}
