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
	"sync"

	"github.com/bogdanovich/dns_resolver"
	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
)

var (
	domain, ip, resUrl, logLevel string
	aDomain                      []string
	cache                        *KvDbOp = NewKvDbOp()
)

type handler struct{}

// 测试在运行的域名范围内
func testIs(s1 string) bool {
	logrus.Info(s1)
	for _, s := range aDomain {
		index := strings.Index(s1, s+".")
		if 0 <= index {
			return true
		}
	}
	return false
}

// 解决相同多次请求的问题
func sendReq(addressOfRequester net.Addr, domain1 string) {
	// 处理过就直接返回，减少 Elasticsearch 服务器交互
	cv, err := cache.Get(domain1)
	if nil != err && "" != string(cv) {
		cache.Put(domain1, []byte(addressOfRequester.String()))
		return
	}
	cache.Put(domain1, []byte(addressOfRequester.String()))
	ip1 := fmt.Sprintf("%s", addressOfRequester)
	a := regexp.MustCompile(`[:]`)
	ip1 = a.Split(ip1, -1)[0]
	i := strings.Count(domain1, "") - 2
	domain1 = domain1[0:i]
	logrus.Info(domain1 + " " + ip1)
	post_body := bytes.NewReader([]byte(fmt.Sprintf(`{"ip":"%s","domain":"%s"}`, ip1, domain1)))
	//
	req, err := http.NewRequest("POST", resUrl, post_body)
	if err == nil {
		// 取消全局复用连接
		// tr := http.Transport{DisableKeepAlives: true}
		// client := http.Client{Transport: &tr}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.2 Safari/605.1.15")
		req.Header.Add("Content-Type", "application/json;charset=UTF-8")
		req.Header.Add("Cache-Control", "no-cache")
		// keep-alive
		req.Header.Add("Connection", "close")
		req.Close = true

		resp, err := http.DefaultClient.Do(req)
		if resp != nil {
			defer resp.Body.Close() // resp 可能为 nil，不能读取 Body
		}
		if err != nil {
			// fmt.Println(err)
			return
		}

		// body, err := ioutil.ReadAll(resp.Body)
		// _, err = io.Copy(ioutil.Discard, resp.Body) // 手动丢弃读取完毕的数据
		// json.NewDecoder(resp.Body).Decode(&data)
		logrus.Info("[send request] " + ip1 + " " + domain1)
		// req.Body.Close()
	}
	// go http.Post(resUrl, "application/json",, post_body)
}

func otherDns(s string) string {
	cv, err := cache.Get(s)
	if nil != err && "" != string(cv) {
		return string(cv)
	}

	resolver := dns_resolver.New([]string{"8.8.8.8", "8.8.4.4"})
	resolver.RetryTimes = 5

	ip, err := resolver.LookupHost(s[0 : strings.Count(s, "")-2])
	if err != nil {
		logrus.Error(err)

	}
	if 0 < len(ip) {
		s1 := fmt.Sprintf(`%s`, ip[0])
		cache.Put(s, []byte(s1))
		return s1
	} else {
		return ""
	}
}

var key, httpHost string
var dnsKm sync.Map

func parseQuery(m *dns.Msg, addressOfRequester net.Addr) {
	for _, q := range m.Question {
		switch q.Qtype {
		// https://golang.hotexamples.com/examples/github.com.miekg.dns/A/Txt/golang-a-txt-method-examples.html
		// https://github.com/joohoi/acme-dns/blob/a33c09accf392a1f2e08628c4766074fc1d1abc1/dns.go#L220
		// answerOwnChallenge
		case dns.TypeTXT:
			{
				fmt.Println(q.Name)
				value1, ok := dnsKm.Load(strings.ToLower(q.Name))
				if !ok {
					fmt.Println(q.Name, "dnsKm.Load=", value1)
					continue
				}
				value := value1.(string)
				// 数字开头的域名不能作为环境变量
				if "" == value {
					fmt.Println("没有对", value)
					value = q.Name
				}
				record := new(dns.TXT)
				record.Hdr = dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeTXT,
					Class:  dns.ClassINET,
					Ttl:    0,
				}
				record.Txt = []string{value}
				m.Answer = append(m.Answer, record)
				continue
			}
		case dns.TypeA, dns.TypeAAAA:
			{
				if testIs(q.Name) {
					// logrus.Info("Query for %s %v\n", q.Name, addressOfRequester)
					rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
					if err == nil {
						m.Answer = append(m.Answer, rr)
					}
				} else {
					szIp1 := otherDns(q.Name)
					if 0 < len(szIp1) {
						logrus.Info("[", szIp1, "]")
						rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, szIp1))
						if err == nil {
							m.Answer = append(m.Answer, rr)
						}
					}
				}
				go sendReq(addressOfRequester, q.Name)
			}
		}
	}
}

func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	addressOfRequester := w.RemoteAddr()
	msg := new(dns.Msg)
	msg.SetReply(r)
	switch r.Opcode {
	case dns.OpcodeQuery, dns.OpcodeIQuery:
		parseQuery(msg, addressOfRequester)
	}

	w.WriteMsg(msg)
}

func HttpApiServer() {
	if "" != httpHost {
		http.HandleFunc("/ACME", func(w http.ResponseWriter, req *http.Request) {
			key1 := req.FormValue("key")
			//log.Println("key1=", key1, "key=", key)
			if key1 == key {
				szDName := req.FormValue("k")
				szDNV := req.FormValue("v")
				//log.Println(szDName, szDNV)
				dnsKm.Store(szDName, szDNV)
				w.Write([]byte("Ok"))
			}
		})
		http.ListenAndServe(httpHost, nil)
	}
}

func main() {
	flag.StringVar(&httpHost, "httpHost", "127.0.0.1:55555", "set ACME http Server ip:port,handle ACME DNS challenges easily,default: 127.0.0.1:55555")
	flag.StringVar(&key, "key", "", "use ACME http API Key")
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

	// http ACME API server
	go HttpApiServer()

	err := srv.ListenAndServe()
	defer srv.Shutdown()
	if err != nil {
		logrus.Fatalf("Failed to set udp listener %s\n", err.Error())
	}

}
