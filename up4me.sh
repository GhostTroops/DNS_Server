#make -f Makefile.cross-compiles
CC=/usr/local/bin/x86_64-linux-musl-gcc CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags=" -s -w " -trimpath -o release/DNS_Server_linux_amd64
scp -i ~/.ssh/id_rsa -C -r -P $newSshPort release/DNS_Server_linux_amd64 root@${newIp}:/root/tools/DNS_Server_linux_amd64
# scp -i ~/.ssh/id_rsa -C -r -P $newSshPort DNS_Server_linux  root@${newIp}:/root/tools/DNS_Server_linux_amd64

