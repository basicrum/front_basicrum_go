#!/bin/bash
STATUSCODE=$(curl --silent --output /dev/stderr --write-out "%{http_code}" http://localhost:14000)
if test $STATUSCODE -ne 400; then
    exit 1
fi