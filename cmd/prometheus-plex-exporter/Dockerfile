FROM alpine:3.24@sha256:3f4d6b32e26703c6082839234f43769d5046b3ac7c910033911d1e71d8cd630a as certs
RUN apk --update --no-cache add ca-certificates tzdata
ARG TARGETARCH
COPY bin/linux/prometheus-plex-exporter-${TARGETARCH} /prometheus-plex-exporter
ENTRYPOINT ["/prometheus-plex-exporter"]