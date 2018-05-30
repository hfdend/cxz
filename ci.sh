#!/bin/bash
GOOS=linux  GOARCH=amd64  go build -ldflags "-s -w"
scp cxz git-sync@47.96.13.222:~/bin/cxz.t
ssh git-sync@cxz /home/git-sync/bin/restart.sh