package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bogdanovich/dns_resolver"
	"github.com/gin-gonic/gin"
	db "github.com/hktalent/goSqlite_gorm/pkg/common"
	db1 "github.com/hktalent/goSqlite_gorm/pkg/db"
	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 记录本地解析
type Result struct {
	gorm.Model
	Dns    string `json:"domain" gorm:"unique_index"`
	Ips    []Ips  `json:"ips" gorm:"many2many:result_ips"`
	Date   string `json:"date"`
	SaveEs bool   `json:"saveEs"`
}

type Ips struct {
	gorm.Model
	Ip string `json:"ip" gorm:"unique_index"`
}

var (
	domain, ip, resUrl, logLevel string
	aDomain                      []string
	cache                        *db.KvDbOp = db.NewKvDbOp()
	dbs                                     = db1.GetDb(&Ips{}, "db/mydbfile")
	doSaveEs                     bool
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
// 一个域名多个ip到情况没有处理
// 在还没有迁移到golang版server该功能先关闭
func sendReq(addressOfRequester net.Addr, domain1 string) {
	if !doSaveEs {
		return
	}
	ip := fmt.Sprintf("%s", addressOfRequester)
	ip1 := ip
	// 跳过不需要记录的dns
	r, err := regexp.Compile(`^(www|ns)\..*\.(51pwn|exploit-poc)\.com`)
	if nil == err {
		a9 := r.FindAllString(strings.ToLower(domain1), -1)
		if nil != a9 && 0 < len(a9) {
			return
		}
	}
	// 处理过就直接返回，减少 Elasticsearch 服务器交互,取到缓存，表示发送过请求
	rD := GetDomain(domain1)
	if nil == rD || !rD.SaveEs {
		a := regexp.MustCompile(`[:]`)
		ip1 = a.Split(ip1, -1)[0]
		i := strings.Count(domain1, "") - 2
		domain1 = domain1[0:i]
		logrus.Info(domain1 + " " + ip1)
		post_body := bytes.NewReader([]byte(fmt.Sprintf(`{"ip":"%s","domain":"%s"}`, ip1, domain1)))
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
			send2cache4R(rD, ip, domain, true)
			// req.Body.Close()
		}
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

func fixdomain(domain1 string) string {
	n1 := len(domain1) - 1
	if "." == domain1[n1:] {
		domain1 = domain1[:n1]
	}
	return domain1
}
func getDate() string {
	var currentTime = time.Now()
	l, err := time.LoadLocation("Asia/Shanghai")
	if nil == err {
		currentTime = time.Now().In(l)
	}
	return currentTime.Format("2006-01-02 15:04:05")
}
func NewResult(ip string, domain1 string, bSave bool) *Result {
	domain1 = fixdomain(domain1)
	var r = &Result{Ips: []Ips{Ips{Ip: ip}}, Dns: domain, SaveEs: bSave}
	r.Date = getDate()
	return r
}
func Result2Byte(r *Result) []byte {
	b, err := json.Marshal(r)
	if nil == err {
		return b
	}
	return nil
}
func Byte2Result(data []byte) *Result {
	var r Result
	err := json.Unmarshal(data, &r)
	if nil != err {
		return nil
	}
	return &r
}

func send2cache4R(r *Result, ip, domain1 string, bSave bool) {
	domain1 = fixdomain(domain1)
	if nil == r {
		r = NewResult(ip, domain1, bSave)
	}
	cache.Put(domain1, []byte(Result2Byte(r)))
}

// 基于缓存记录dns解析日志
func send2cache(addressOfRequester net.Addr, domain1 string, bSave bool) {
	domain1 = fixdomain(domain1)
	cv, err := cache.Get(domain1)
	s1 := addressOfRequester.String()
	if nil == err && s1 == string(cv) {
		return
	}
	cache.Put(domain1, []byte(Result2Byte(NewResult(ip, domain1, bSave))))
}

func GetDomain(d string) *Result {
	d = fixdomain(d)
	cv, err := cache.Get(d)
	if nil == err {
		r := Byte2Result(cv)
		return r
	}
	return nil
}

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
					logrus.Debug(q.Name, " dnsKm.Load= ", value1)
					continue
				}
				value := value1.(string)
				// 数字开头的域名不能作为环境变量
				if "" == value {
					logrus.Debug("没有对", value)
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
				go send2cache(addressOfRequester, q.Name, false)
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

func dnsRes(g *gin.Context) {
	x1 := dbs
	if nil == x1 {
		logrus.Debug("dbs is nil,db1.GetDb 失败了")
		x1 = db1.GetDb(&Ips{}, "db/mydbfile")
		if nil != x1 {
			dbs = x1
		} else {
			logrus.Debug("db1.GetDb 失败了")
		}
	}
	req := g.Request
	q := req.FormValue("q")
	q = strings.TrimSpace(q)
	if "" != q {
		r := GetDomain(q)
		if nil != r {
			g.JSON(http.StatusOK, r)
			return
		} else if nil != x1 {
			var rst []Result = db1.GetSubQueryLists(Result{}, "Ips",
				[]Result{}, 10, 0, "dns = ?", q)
			if 0 < len(rst) {
				if 0 < len(rst[0].Ips) {
					g.JSON(http.StatusOK, map[string]interface{}{"domain": q, "ip": rst[0].Ips[0].Ip, "date": rst[0].Date})
					return
				} else {
					logrus.Debug(rst[0])
					logrus.Debug(rst[0].Ips)
				}
				//var r1 Result
				//r2 := db1.GetOne[Result](&r1, "dns=?", q)
				//if nil != r2 {
				//g.JSON(http.StatusOK, r2)
			} else {
				logrus.Debug("db1.GetOne not found ", q)
			}

		}
	}
	g.JSON(http.StatusNotFound, gin.H{"msg": "not found"})
}
func ACME(g *gin.Context) {
	req := g.Request
	key1 := req.FormValue("key")
	logrus.Debug("req.FormValue key1=", key1, "key=", key)
	if key1 == key {
		szDName := req.FormValue("k")
		szDNV := req.FormValue("v")
		logrus.Debug(szDName, szDNV)
		dnsKm.Store(szDName, szDNV)
		g.JSON(http.StatusOK, "ok")
	}
}

func getip(g *gin.Context) {
	ip, ok := g.Request.Header["X-Real-Ip"]
	//host, ok1 := g.Request.Header["Host"]
	//if ok && ok1 && 0 < len(host) && 0 < len(ip) && strings.HasPrefix(host[0], "ip.") {
	if ok && 0 < len(ip) {
		g.JSON(http.StatusOK, ip[0])
		return
	}
	//logrus.Debug(ip, g.Request.Header)
	g.JSON(http.StatusBadRequest, "can not get ip")
	//return false
}
func ip2domain(g *gin.Context) {
	if nil == dbs {
		logrus.Debug("dbs is nil")
		return
	}
	var r Result
	err := g.BindJSON(&r)
	if nil == err {
		r.Date = getDate()
		if 1 == db1.Create[Result](&r) {
			logrus.Debug("save ok ", r.Dns)
		} else {
			logrus.Debug("not save ok ", r.Dns)
		}
		g.JSON(http.StatusOK, "ok")
		return
	} else {
		logrus.Debug("BindJSON ", err)
	}
	g.JSON(http.StatusBadRequest, err)
}
func HttpApiServer() {
	if "" != httpHost {
		gin.SetMode(gin.ReleaseMode)
		router := gin.Default()
		s1 := "/ACME"
		router.GET(s1, ACME)
		router.POST(s1, ACME)
		s1 = "/ip2domain"
		router.GET(s1, ip2domain)
		router.GET("/", getip)
		router.POST(s1, ip2domain)
		router.GET("/dnslog", dnsRes)
		router.Run(httpHost)
	}
}

func main() {
	var ExpiresAt uint64
	flag.StringVar(&httpHost, "httpHost", "127.0.0.1:55555", "set ACME http Server ip:port,handle ACME DNS challenges easily,default: 127.0.0.1:55555")
	flag.StringVar(&key, "key", "", "use ACME http API Key")
	flag.StringVar(&domain, "domain", "51pwn.com,exploit-poc.com", "set domain eg: 51pwn.com")
	flag.StringVar(&ip, "ip", "144.34.164.150", "set domain server ip, eg: 222.44.11.3")
	flag.StringVar(&resUrl, "resUrl", "", "Set the Elasticsearch url that accepts dns parsing logs, eg: http://127.0.0.1/dnsRecode")
	flag.StringVar(&logLevel, "level", "WARN", "set loglevel, option")
	flag.Uint64Var(&ExpiresAt, "ExpiresAt", 120000, "default 120s = 120000")
	flag.Parse()
	//cache.SetExpiresAt(ExpiresAt)
	dbs.AutoMigrate(&Ips{}, &Result{})
	doSaveEs = "" != resUrl

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
