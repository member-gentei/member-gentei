NOWISH := $(shell date +%s)

.PHONY: bot-rpm
bot-rpm:
	rm -rf gentei-bot-*.rpm build
	mkdir -p build
	install -d -m 755 build/bin
	install -d -m 755 build/etc/gentei
	install -d -m 755 build/etc/systemd/system
	cd bot && \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../build/bin/gentei-bot .
	install -D -m 600 init/bot.yml build/etc/gentei/bot.yml
	install -D -m 644 init/gentei-bot.service build/etc/systemd/system/gentei-bot.service
	fpm -s dir \
		-t rpm \
		-C build \
		--name gentei-bot \
		--version '$(NOWISH)' \
		--maintainer 'Mark Ignacio <mark@ignacio.io>' \
		--before-install init/scripts/bot/before-install.sh \
		--config-files /etc/gentei/ \
		--rpm-user gentei-bot \
		--rpm-group gentei-bot \
		--rpm-use-file-permissions \
		.
	rm -rf build.PHONY: bot-rpm

.PHONY: member-check-rpm
member-check-rpm:
	rm -rf gentei-member-check-*.rpm build
	mkdir -p build
	install -d -m 755 build/bin
	install -d -m 755 build/etc/gentei
	install -d -m 755 build/etc/systemd/system
	cd jobs/member-check && \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../build/bin/gentei-member-check .
	install -D -m 600 init/member-check.yml build/etc/gentei/member-check.yml
	install -D -m 644 init/gentei-member-check.service build/etc/systemd/system/gentei-member-check.service
	install -D -m 644 init/gentei-member-check.timer build/etc/systemd/system/gentei-member-check.timer
	fpm -s dir \
		-t rpm \
		-C build \
		--name gentei-member-check \
		--version '$(NOWISH)' \
		--maintainer 'Mark Ignacio <mark@ignacio.io>' \
		--before-install init/scripts/member-check/before-install.sh \
		--config-files /etc/gentei/ \
		--rpm-user gentei-member-check \
		--rpm-group gentei-member-check \
		--rpm-use-file-permissions \
		.
	rm -rf build

.PHONY: refresh-data-rpm
refresh-data-rpm:
	rm -rf gentei-refresh-data-*.rpm build
	mkdir -p build
	install -d -m 755 build/bin
	install -d -m 755 build/etc/gentei
	install -d -m 755 build/etc/systemd/system
	cd jobs/refresh-data && \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../build/bin/gentei-refresh-data .
	install -D -m 600 init/refresh-data.yml build/etc/gentei/refresh-data.yml
	install -D -m 644 init/gentei-refresh-data.service build/etc/systemd/system/gentei-refresh-data.service
	install -D -m 644 init/gentei-refresh-data.timer build/etc/systemd/system/gentei-refresh-data.timer
	fpm -s dir \
		-t rpm \
		-C build \
		--name gentei-refresh-data \
		--version '$(NOWISH)' \
		--maintainer 'Mark Ignacio <mark@ignacio.io>' \
		--before-install init/scripts/refresh-data/before-install.sh \
		--config-files /etc/gentei/ \
		--rpm-user gentei-refresh-data \
		--rpm-group gentei-refresh-data \
		--rpm-use-file-permissions \
		.
	rm -rf build

.PHONY: rpms
rpm: bot-rpm member-check-rpm refresh-data-rpm