#!/bin/bash

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