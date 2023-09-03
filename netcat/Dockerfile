FROM alpine:latest
RUN apk update && apk add netcat-openbsd
COPY netcat/netcat.sh /
RUN chmod +x netcat.sh
ENTRYPOINT ["/netcat.sh"]