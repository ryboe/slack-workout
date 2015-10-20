FROM busybox
MAINTAINER y0ssar1an
COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY slack-pushups /slack-pushups
EXPOSE 8000
ENTRYPOINT ["/slack-pushups"]
