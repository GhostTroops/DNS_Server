[![Tweet](https://img.shields.io/twitter/url/http/Hktalent3135773.svg?style=social)](https://twitter.com/intent/follow?screen_name=Hktalent3135773) [![Follow on Twitter](https://img.shields.io/twitter/follow/Hktalent3135773.svg?style=social&label=Follow)](https://twitter.com/intent/follow?screen_name=Hktalent3135773) [![GitHub Followers](https://img.shields.io/github/followers/hktalent.svg?style=social&label=Follow)](https://github.com/hktalent/)
[![Top Langs](https://profile-counter.glitch.me/hktalent/count.svg)](https://51pwn.com)

# DNS_Server
Simple DNS log Server,free DNS log server

# what's the new
- Reduce requests to Elasticsearch servers based on caching
- Added ACME DNS challenge design

# How Build
```bash
git clone https://github.com/hktalent/DNS_Server
cd DNS_Server
go mod init github.com/facebookincubator/nvdtools
go mod tidy
# go build

make -f Makefile.cross-compiles

```
## how cross build in macos
```bash
brew install FiloSottile/musl-cross/musl-cross
which gobuild
/usr/local/bin/gobuild
cat  /usr/local/bin/gobuild
go build -ldflags="-s -w " -trimpath -o $2 $1
CC=x86_64-linux-musl-gcc CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-linkmode external -extldflags -static -s -w " -trimpath -o $2_linux $1

gobuild . DNS_Server

```

# How use
```bash
git commit -m "fix" .;git push
cp build/bin/* ~/go/bin/

scp -i ~/.ssh/id_rsa -r -P $myVpsPort DNS_Server root@51pwn.com:/root/
```

# How ACME DNS challenge
## run server by key
```
./DNS_Server -key="dd9j-dds-33xfgk-33"
```
vi upApi.sh
```
curl -s -v "http://127.0.0.1:55555/ACME?key=dd9j-dds-33xfgk-33&k=${1}&v=${2}"
```
## challenge
```
docker run -it -p 80:80  --rm --name certbot \
     -v "/etc/letsencrypt:/etc/letsencrypt" \
     -v "/var/lib/letsencrypt:/var/lib/letsencrypt" \
     certbot/certbot certonly  -d *.51pwn.com -d 51pwn.com -d *.exploit-poc.com  -d exploit-poc.com \
     --manual --preferred-challenges dns --server https://acme-v02.api.letsencrypt.org/directory

chmod +x upApi.sh

./upApi.sh '_acme-challenge.exploit-poc.com.'.'QZPE4B9OQivKZDi7Hq3On1IhhdZiEX2iVJ8ojKuOGsA'
./upApi.sh "_acme-challenge.51pwn.com" "q31hmgemyDDsU_rTIM8cW3h0EExs0HPt-SqwoVa0AV8"
```
After running the above command, confirm the certonly operation steps

QZPE4B9OQivKZDi7Hq3On1IhhdZiEX2iVJ8ojKuOGsA and cW3h0EExs0HPt-SqwoVa0AV8 replace your
<img width="832" alt="image" src="https://user-images.githubusercontent.com/18223385/164650694-b35d3e6d-6b8b-4f01-b8c2-9d54269fd53c.png">

<img width="813" alt="image" src="https://user-images.githubusercontent.com/18223385/164650641-cabde872-e7df-4a66-b55e-92a95103e24b.png">
<img width="585" alt="image" src="https://user-images.githubusercontent.com/18223385/164710659-a30b12b8-8059-4489-843e-15d1c773a860.png">
# How config dns
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
# Donation
| Wechat Pay | AliPay | Paypal | BTC Pay |BCH Pay |
| --- | --- | --- | --- | --- |
|<img src=https://github.com/hktalent/myhktools/blob/master/md/wc.png>|<img width=166 src=https://github.com/hktalent/myhktools/blob/master/md/zfb.png>|[paypal](https://www.paypal.me/pwned2019) **miracletalent@gmail.com**|<img width=166 src=https://github.com/hktalent/myhktools/blob/master/md/BTC.png>|<img width=166 src=https://github.com/hktalent/myhktools/blob/master/md/BCH.jpg>|


