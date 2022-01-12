# DNS_Server
DNS_Server

# How Build
```bash
git clone https://github.com/hktalent/DNS_Server
cd DNS_Server
go mod init github.com/facebookincubator/nvdtools
go mod tidy
# go build
go install github.com/karalabe/xgo@latest
xgo .
```
# How use
```bash
git commit -m "fix" .;git push
cp build/bin/* ~/go/bin/

scp -i ~/.ssh/id_rsa -r -P $myVpsPort DNS_Server root@51pwn.com:/root/
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

