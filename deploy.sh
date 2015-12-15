#!/usr/bin/env bash
set -e  # exit immediately if any command fails

GOOS=linux go build -ldflags '-w' slack-workout.go

set +e
docker rmi y0ssar1an/slack-workout

set -e
docker build -t y0ssar1an/slack-workout .
docker push y0ssar1an/slack-workout

ssh -c "sudo systemctl restart slack-workout" coreos-01
