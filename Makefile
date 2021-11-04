COMMIT := $(shell git describe --tags --always --dirty)
VERSION := dev

OS_ARCH :=$(shell go env GOHOSTOS)_$(shell go env GOHOSTARCH)
BIN := terraform-provider-bunny
INSTALLDIR := "$(HOME)/.terraform.d/plugins/local/simplesurance/bunny/"

LDFLAGS := "-X github.com/simplesurance/terraform-provider-bunny/internal/provider.Version=$(VERSION) -X github.com/simplesurance/terraform-provider-bunny/internal/provider.Commit=$(COMMIT)"
BUILDFLAGS := -trimpath -ldflags=$(LDFLAGS)

TFTRC_FILENAME := bunny-dev.tftrc

SWEEPERS := pullzones

REPO_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

default: build

.PHONY: build
build:
	$(info * compiling $(BIN))
	go build $(BUILDFLAGS) -o $(BIN)

.PHONY: check
check:
	$(info * running golangci-lint code checks)
	golangci-lint run

.PHONY: install
install:
	$(info * installing provider to $(INSTALLDIR)/$(BIN))
	GOBIN=$(INSTALLDIR) go install $(BUILDFLAGS)

.PHONY: gen-dev-tftrc
gen-dev-tftrc:
	$(info * generating $(TFTRC_FILENAME))
	@scripts/gen-dev-tftrc.sh "$(INSTALLDIR)" > $(TFTRC_FILENAME)
	$(info run 'export TF_CLI_CONFIG_FILE=$(REPO_ROOT)/$(TFTRC_FILENAME)' to use it)


# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test $(BUILDFLAGS) -race -v $(TESTARGS) -timeout 120m ./...

.PHONY: docs
docs:
	$(info * generating documentation files (docs/))
	go generate $(BUILDFLAGS)

.PHONY: sweep
sweep:
	$(warning WARNING: This will destroy infrastructure. Use only in development accounts.)
	cd internal/provider && TF_ACC=1 go test -v -sweep=all -sweep-run=$(SWEEPERS) -timeout=5m ./...
