#!/bin/bash

/usr/local/nginx/sbin/nginx -c /usr/local/nginx/conf/nginx.conf
cd /root/go/src/github.com/arthurkiller/shadowsocks-go/cmd/shadowsocks-server
./shadowsocks-server -config=config.json 1>1.txt 2>2.txt &
cd /root/go/src/github.com/dilfish/libsm
/usr/bin/nohup ./libsm > libsm.log & 
cd /disk1/cao/91
/usr/bin/nohup ./dl91go > dl91go.log &
