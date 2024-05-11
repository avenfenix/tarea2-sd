#!/bin/bash
# MAQUINA VIRTUAL 3

# GOLANG
cd ~/
wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz
rm -rf go1.22.3.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# RABBITMQ SERVER
sudo apt-get update
sudo apt-get -y install socat logrotate init-system-helpers adduser
sudo apt-get -y install rabbitmq-server
sudo rabbitmqctl add_user admin 1234
sudo rabbitmqctl set_user_tags admin administrator
sudo rabbitmqctl set_permissions -p / admin ".*" ".*" ".*"
sudo systemctl restart rabbitmq-server