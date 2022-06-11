make -f Makefile.cross-compiles
scp -i ~/.ssh/id_rsa -C -r -P $newSshPort release/DNS_Server_linux_amd64 root@${newIp}:/root/tools/DNS_Server_linux_amd64
# scp -i ~/.ssh/id_rsa -C -r -P $newSshPort DNS_Server_linux  root@${newIp}:/root/tools/DNS_Server_linux_amd64

