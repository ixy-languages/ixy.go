#!/bin/bash
V=$1
if [ -z "$0" ]; then
	V=1.12.5
fi
wget -P ~/ https://dl.google.com/go/go$V.linux-amd64.tar.gz
tar -C /usr/local -xzf ~/go$V.linux-amd64.tar.gz
rm ~/go$V.linux-amd64.tar.gz
echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.profile
source ~/.profile
mkdir ~/go/
mkdir ~/go/src/
#cd ~/go/src/
#git clone https://github.com/ixy-languages/ixy.go.git
