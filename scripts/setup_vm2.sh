#!/bin/bash
# MAQUINA VIRTUAL 2

# GOLANG
cd ~/
wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz
rm -rf go1.22.3.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# MONGODB
sudo mkdir /data/
sudo mkdir /data/db/
sudo apt update
sudo apt install -y gnupg wget apt-transport-https ca-certificates software-properties-common
wget -qO- \
  https://pgp.mongodb.com/server-7.0.asc | \
  gpg --dearmor | \
  sudo tee /usr/share/keyrings/mongodb-server-7.0.gpg >/dev/null

echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-7.0.gpg ] \
  https://repo.mongodb.org/apt/ubuntu $(lsb_release -cs)/mongodb-org/7.0 multiverse" | \
  sudo tee -a /etc/apt/sources.list.d/mongodb-org-7.0.list
sudo apt update
sudo apt install -y mongodb-org
sudo systemctl restart mongod