FROM alpine:3.24@sha256:71e9b7bae160ba167fa17414401b94cbcdaf3c32feeab9298bd794a2a0614c35 as certs
RUN apk --update --no-cache add ca-certificates tzdata
ARG TARGETARCH
COPY bin/linux/prometheus-plex-exporter-${TARGETARCH} /prometheus-plex-exporter
ENTRYPOINT ["/prometheus-plex-exporter"]