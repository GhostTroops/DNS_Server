# DNS_Server
DNS_Server

# How Build
```bash
go mod init github.com/facebookincubator/nvdtools
go mod tidy
go build
git commit -m "fix" .;git push
cp build/bin/* ~/go/bin/

scp -i ~/.ssh/id_rsa -r -P $myVpsPort DNS_Server root@51pwn.com:/root/
```

# How Row
```bash
./DNS_Server -ip 199.180.115.7 -domain 51pwn.com
```

