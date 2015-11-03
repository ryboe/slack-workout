#!/usr/bin/env bash

GOOS=linux go build -ldflags '-w' slack-pushups.go

docker rmi y0ssar1an/slack-pushups
docker build -t y0ssar1an/slack-pushups .
docker push y0ssar1an/slack-pushups

ssh -c "sudo systemctl restart slack-pushups" coreos-01
