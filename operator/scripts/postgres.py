#! /usr/bin/env python
import pexpect
import sys
import os
import re
from pexpect import pxssh

def remove_ansi_escape_sequences(string):
    ansi_escape = re.compile(r'\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])')
    return ansi_escape.sub('', string)

def remove_cr(string):
    r = re.compile(r'\r')
    return r.sub('', string)

hostname=os.environ['HOSTNAME']
username=os.environ['USERNAME']
password=os.environ['PASSWORD']

s = pxssh.pxssh(options={"StrictHostKeyChecking": "no", "UserKnownHostsFile": "/dev/null"})
s.force_password = True
s.login (hostname, username, password)

# If any rows return 'Y' then we are good.
s.sendline('''echo "select count(*) from data where data_col = 'Y'" | psql -t argus_demo | xargs''')
s.prompt()
output=s.before.decode('utf-8').split('\n')
line=remove_ansi_escape_sequences(output[1].strip())
line=remove_cr(line)
count=int(line)

if count > 0:
    sys.exit(0)
else:
    sys.exit(1)
