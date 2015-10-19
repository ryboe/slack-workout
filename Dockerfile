FROM busybox
MAINTAINER y0ssar1an
COPY slack-pushups /slack-pushups
ENV SLACK_API_TOKEN $SLACK_API_TOKEN
EXPOSE 8000
ENTRYPOINT ["/slack-pushups"]
