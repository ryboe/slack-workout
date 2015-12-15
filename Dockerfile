FROM busybox
MAINTAINER y0ssar1an
COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY zoneinfo /usr/share/zoneinfo
COPY slack-workout /slack-workout
EXPOSE 8000
ENTRYPOINT ["/slack-workout"]
