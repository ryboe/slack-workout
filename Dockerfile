FROM busybox
MAINTAINER y0ssar1an
COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY zoneinfo /usr/share/zoneinfo
COPY slack-workout /slack-workout
EXPOSE 8001
ENTRYPOINT ["/slack-workout"]
