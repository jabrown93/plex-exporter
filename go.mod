module github.com/grafana/plexporter

go 1.25.0

require (
	github.com/go-kit/log v0.2.1
	github.com/gorilla/websocket v1.5.3
	github.com/jrudio/go-plex-client v0.0.0-20250127195314-943dc7a39f7c
	github.com/prometheus/client_golang v1.24.1
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.70.1 // indirect
	github.com/prometheus/procfs v0.21.1 // indirect
	golang.org/x/sys v0.47.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

// Replace for fix: https://github.com/jrudio/go-plex-client/pull/56
replace github.com/jrudio/go-plex-client => github.com/jsclayton/go-plex-client v0.0.0-20230428232959-d53064b6f34a
