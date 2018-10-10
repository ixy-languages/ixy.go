#!/bin/bash
wget -P ~/ https://dl.google.com/go/go1.11.1.linux-amd64.tar.gz
tar -C /usr/local -xzf ~/go1.11.1.linux-amd64.tar.gz
rm ~/go1.11.1.linux-amd64.tar.gz
echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.profile
source ~/.profile
mkdir ~/go/
mkdir ~/go/src/
cd ~/go/src/
git clone https://github.com/ixy-languages/ixy.go.git
