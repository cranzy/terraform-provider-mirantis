TERRAFORM_PROVIDER_ROOT=mirantis.com/providers
BINARY_ROOT=terraform-provider
INSTALL_ROOT?=$(HOME)/.terraform.d/plugins
LOCAL_BIN_PATH?=./bin

VERSION?=$(shell git describe --tags)

PROVIDERS?=msr
ARCHES?=amd64 arm64
OSES?=linux darwin

GO=$(shell which go)

default: install

clean:
	rm -rf "$(LOCAL_BIN_PATH)"
	rm -rf "$(INSTALL_ROOT)/$(TERRAFORM_PROVIDER_ROOT)"

build:
	mkdir -p $(LOCAL_BIN_PATH)
	for PROVIDER in $(PROVIDERS); do \
		for OS in $(OSES); do \
			for ARCH in $(ARCHES); do \
				GOOS=$$OS GOARCH=$$ARCH $(GO) build -v -o "$(LOCAL_BIN_PATH)/$(BINARY_ROOT)-$$PROVIDER-$$OS_$$ARCH" "./cmd/$$PROVIDER"; \
			done; \
		done; \
	done;

release:
	goreleaser release --rm-dist --snapshot --skip-publish  --skip-sign

install: build
	for PROVIDER in $(PROVIDERS); do \
		for OS in $(OSES); do \
			for ARCH in $(ARCHES); do \
				mkdir -p "$(INSTALL_ROOT)/$(TERRAFORM_PROVIDER_ROOT)/$$PROVIDER/$(VERSION)/$${OS}_$${ARCH}"; \
				cp "$(LOCAL_BIN_PATH)/$(BINARY_ROOT)-$$PROVIDER-$$OS_$$ARCH" "$(INSTALL_ROOT)/$(TERRAFORM_PROVIDER_ROOT)/$$PROVIDER/$(VERSION)/$${OS}_$${ARCH}/$(BINARY_ROOT)-$$PROVIDER"; \
    	done; \
		done; \
	done;


test:
	go test -i ./...

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: clean build install test testacc