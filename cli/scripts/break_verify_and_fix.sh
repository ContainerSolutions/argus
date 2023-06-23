#! /bin/bash

fail=0
az vm stop -g vm-lb-test -n $1
#sleeping 60 seconds for vm to be turned off
sleep 15
curl --connect-timeout 5 -w '%{http_code}\n' 'http://20.8.65.96/' || fail=1
az vm start -g vm-lb-test -n $1
exit $fail