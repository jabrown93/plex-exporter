FROM alpine:3.24@sha256:be2656061b536c1883c4034790e21d0a09c5aaaf51968b89368a3c0371ea2c9a as certs
RUN apk --update --no-cache add ca-certificates tzdata
ARG TARGETARCH
COPY bin/linux/prometheus-plex-exporter-${TARGETARCH} /prometheus-plex-exporter
ENTRYPOINT ["/prometheus-plex-exporter"]