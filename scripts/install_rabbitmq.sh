#!/bin/bash

# RABBITMQ SERVER
sudo apt-get update
sudo apt-get -y install socat logrotate init-system-helpers adduser
sudo apt-get -y install rabbitmq-server
sudo rabbitmqctl add_user admin 1234
sudo rabbitmqctl set_user_tags admin administrator
sudo rabbitmqctl set_permissions -p / admin ".*" ".*" ".*"
sudo systemctl restart rabbitmq-server