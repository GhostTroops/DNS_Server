# DNS_Server
Simple DNS log Server,free DNS log server

# How Build
```bash
git clone https://github.com/hktalent/DNS_Server
cd DNS_Server
go mod init github.com/facebookincubator/nvdtools
go mod tidy
# go build
go install github.com/karalabe/xgo@latest
xgo .

# or 
make -f Makefile.cross-compiles

```
# How use
```bash
git commit -m "fix" .;git push
cp build/bin/* ~/go/bin/

scp -i ~/.ssh/id_rsa -r -P $myVpsPort DNS_Server root@51pwn.com:/root/
```

#How config dns
- add  hosts
```
https://dcc.godaddy.com/manage/51pwn.com/dns/hosts
ns1.51pwn.com
ns2.51pwn.com

https://dcc.godaddy.com/manage/exploit-poc.com/dns/hosts
ns1.exploit-poc.com
ns2.exploit-poc.com
```
- change dns to
```
https://dcc.godaddy.com/manage/51pwn.com/dns
ns1.51pwn.com
ns2.51pwn.com

https://dcc.godaddy.com/manage/exploit-poc.com/dns
ns1.exploit-poc.com
ns2.exploit-poc.com
```

# How Run
```bash
./DNS_Server -ip 199.180.115.7 -domain 51pwn.com,exploit-poc.com -resUrl http://127.0.0.1:9999

```
now,you can test
```bash
ping sdfsfs.51pwn.com
ping sdlfdslfjdslkfj.exploit-poc.com

```
# Example
```
curl -v 'ldap://xx22.log4j2.51pwn.com'
# or 
ping  -c 2 'xx22.log4j2.51pwn.com'

# how get check reuslt
curl -v -H 'user-agent: Mozilla/5.0 (Windows NT 6.1; rv:45.0) Gecko/20100101 Firefox/45.0' -k -o- 'https://51pwn.com/dnslog?q=xx22.log4j2.51pwn.com'
```
### result
```json
{"ip":"172.70.209.200","domain":"xx22.log4j2.51pwn.com","type":"dnslog","date":"2022-01-23 05:19:53"}                                                                                             
```

