#! /usr/bin/env python
import pexpect
import sys
import os
from pexpect import pxssh

hostname=os.environ['HOSTNAME']
username=os.environ['USERNAME']
password=os.environ['PASSWORD']

s = pxssh.pxssh(options={"StrictHostKeyChecking": "no", "UserKnownHostsFile": "/dev/null"})
s.force_password = True
s.login (hostname, username, password)

s.sendline('echo blah')
s.prompt()
print(s.before.decode('utf-8').split('\n', 1)[1])
# If output matches string, sys.exit(0), if it doesn't, sys.exit(1)
