FROM alpine:3.6

RUN apk --no-cache add ca-certificates tini curl bash


EXPOSE 9099

COPY bin/prometheus_webhook_logger /prometheus_webhook_logger
COPY sample-alert.json /

ENTRYPOINT ["/bin/bash", "-c", "/prometheus_webhook_logger \"$@\"", "--"]