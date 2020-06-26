#!/usr/bin/env python3
#sudo systemctl start redis-server
#sudo systemctl stop redis-sentinel
#sudo systemctl stop redis-server
import socket
import time
import redis

r = redis.Redis(host='localhost', port=8090)
for x in range(1):
    a = r.execute_command('set foo 2')
    print(a)

for x in range(1):
    a = r.execute_command('get foo')
    print(a)

r.close()
