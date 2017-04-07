FROM busybox:latest
MAINTAINER Michael Dai <sarahdj0917@gmail.com>

COPY xunyu /bin/xunyu
COPY config.json /etc/xunyu/config.json

VOLUME ["/xunyu"]
WORKDIR /xunyu
ENTRYPOINT ["/bin/xunyu"]
CMD ["-config=/etc/xunyu/config.json"]