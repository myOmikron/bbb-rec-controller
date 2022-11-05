SHELL = /usr/bin/env bash

# Build directory
BUILD_DIR = bin

# Destination directory
DEST_DIR = /usr/local/bin

.PHONY: build clean install uninstall

build: clean
	go build -ldflags=-w -o ${BUILD_DIR}/ ./...

clean:
	rm -rf ${BUILD_DIR}/

install:
	systemctl stop bbb-rec-controller.service ||:
	mkdir -p /etc/bbb-rec-controller/
	cp example.config.toml /etc/bbb-rec-controller/
	cp -r ${BUILD_DIR}/bbb-rec-controller ${DEST_DIR}/
	cp bbb-rec-controller.service /lib/systemd/system/
	if [ -L /etc/systemd/system/multi-user.target.wants/bbb-rec-controller.service ] ; then \
		if [ -e /etc/systemd/system/multi-user.target.wants/bbb-rec-controller.service ]; then \
			echo "Service file is already linked properly"; \
		else \
			rm /etc/systemd/system/multi-user.target.wants/bbb-rec-controller.service; \
			ln -s /lib/systemd/system/bbb-rec-controller.service /etc/systemd/system/multi-user.target.wants/; \
		fi \
	else \
		ln -s /lib/systemd/system/bbb-rec-controller.service /etc/systemd/system/multi-user.target.wants/; \
	fi
	systemctl daemon-reload
	systemctl enable bbb-rec-controller

uninstall:
	rm ${DEST_DIR}/bbb-rec-controller
	rm -rf /etc/bbb-rec-controller/ /var/lib/bbb-rec-controller/