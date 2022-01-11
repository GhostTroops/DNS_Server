# DNS_Server
DNS_Server

# How Build
```bash
go mod init github.com/facebookincubator/nvdtools
go mod tidy
go build
cp build/bin/* ~/go/bin/

scp -i ~/.ssh/id_rsa -r -P $myVpsPort DNS_Server root@51pwn.com:/root/
```
