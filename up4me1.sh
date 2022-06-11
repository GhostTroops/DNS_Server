rm -rf release
rm -rf DNS_Server_linux DNS_Server
go mod vendor
go mod verify
gobuild . DNS_Server
scp -i ~/.ssh/id_rsa -C -r -P $newSshPort DNS_Server_linux  root@${newIp}:/root/tools/DNS_Server_linux_amd64
rm -rf DNS_Server_linux
rm -rf DNS_Server
