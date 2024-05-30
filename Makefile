SHELL=bash

include .env

dev:
	$(eval export $(sed 's/#.*//g' .env | xargs))
	(cd cmd/sentry && go run main.go)

test:
	go test -json -skip /pkg/test -v ./... $(args) 2>&1 | gotestfmt

test-match:
	make test args="-run $(case)"

test-ci:
	set -euo pipefail
	go test -json -skip /pkg/test -v ./... 2>&1 | gotestfmt

vulns-check:
	govulncheck ./...

docs:
	pkgsite -http=:4060