
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
	cd gentei && yarn build
	firebase deploy

.PHONY:
i18n:
	cd bot && \
		goi18n extract -sourceLanguage en-US -outdir i18n && \
		goi18n merge -sourceLanguage en-US -outdir i18n i18n/active.en-US.toml i18n/active.zh-Hant.toml && \
		goi18n merge -sourceLanguage en-US -outdir i18n i18n/active.*.toml i18n/translate.*.toml
	rm -f bot/i18n/translate.*.toml

.PHONY: translations
translations:
	cd tools && go run main.go botgen

.PHONY: clean
clean:
	rm -rf *.rpm build
