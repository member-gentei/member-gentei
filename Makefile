
include *.mk

.PHONY: openapi
openapi:
	oapi-codegen -generate 'types,chi-server,spec' -package api openapi.yml > functions/API/api.gen.go
	oapi-codegen -generate 'types,client,spec' -package api openapi.yml > bot/discord/api/api.gen.go

.PHONY: docs
docs:
	mv docs/CNAME public-docs/
	rm -rf docs
	cd public-docs && \
		doctave build --release
	mv public-docs/CNAME public-docs/site/
	mv public-docs/site docs
	yarn run redoc-cli bundle openapi.yml --cdn -o docs/api.html

.PHONY: api
api: docs function-api

.PHONY: firebase
firebase:
	firebase deploy

.PHONY: clean
clean:
	rm -rf *.rpm build
