#!/bin/bash
# MAQUINA VIRTUAL 1

cd ~/
wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz
rm -rf go1.22.3.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc