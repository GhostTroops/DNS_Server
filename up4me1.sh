rm -rf release
go mod vendor
go mod verify
gobuild . DNS_Server
scp -i ~/.ssh/id_rsa -C -r -P $newSshPort release/DNS_Server_linux  root@${newIp}:/root/tools/DNS_Server_linux_amd64
rm -rf release
