COMMIT := $(shell git describe --tags --always --dirty)
VERSION := 0.0.1

OS_ARCH :=$(shell go env GOHOSTOS)_$(shell go env GOHOSTARCH)
BIN := terraform-provider-bunny
INSTALLDIR := "$(HOME)/.terraform.d/plugins/registry.terraform.io/simplesurance/bunny/$(VERSION)/$(OS_ARCH)"

LDFLAGS := "-X github.com/simplesurance/terraform-provider-bunny/internal/provider.Version=$(VERSION) -X github.com/simplesurance/terraform-provider-bunny/internal/provider.Commit=$(COMMIT)"
BUILDFLAGS := -trimpath -ldflags=$(LDFLAGS)

SWEEPERS := pullzones

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
	GOBIN=$(INSTALLDIR) go install $(BUILDFLAGS)
	@echo
	@echo Provider installed to $(INSTALLDIR)/$(BIN)

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test $(BUILDFLAGS) -race -v $(TESTARGS) -timeout 120m ./...

.PHONY: docs
docs:
	go generate $(BUILDFLAGS)

.PHONY: sweep
sweep:
	$(warning WARNING: This will destroy infrastructure. Use only in development accounts.)
	cd internal/provider && TF_ACC=1 go test -v -sweep=all -sweep-run=$(SWEEPERS) -timeout=5m ./...
