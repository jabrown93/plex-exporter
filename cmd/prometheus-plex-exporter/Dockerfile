FROM alpine:3.24@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6 as certs
RUN apk --update --no-cache add ca-certificates tzdata
ARG TARGETARCH
COPY bin/linux/prometheus-plex-exporter-${TARGETARCH} /prometheus-plex-exporter
ENTRYPOINT ["/prometheus-plex-exporter"]