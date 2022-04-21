rm -rf release
scp -i ~/.ssh/id_rsa -C -r -P $newSshPort release/DNS_Server_linux_amd64 root@${newIp}:/root/

