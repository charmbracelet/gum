.PHONY: build test lint clean install uninstall run panel

BINARY_NAME=gum
INSTALL_PATH=$(DESTDIR)/usr/local/bin
GO=go
GOFLAGS=-v

build:
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) .

test:
	$(GO) test -v ./panel/

test-all:
	$(GO) test -v ./...

lint:
	golangci-lint run ./panel/

lint-all:
	golangci-lint run ./...

clean:
	rm -f $(BINARY_NAME)

install: build
	install -m 755 $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)

uninstall:
	rm -f $(INSTALL_PATH)/$(BINARY_NAME)

run-panel:
	$(GO) run . panel choose a b c filter x y z

run-filter:
	$(GO) run . filter apple banana cherry

run-choose:
	$(GO) run . choose apple banana cherry

build-all:
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BINARY_NAME)-windows-amd64.exe .

help:
	@echo "Available targets:"
	@echo "  build        - Build gum binary"
	@echo "  test         - Run panel tests"
	@echo "  test-all     - Run all tests"
	@echo "  lint         - Run linter on panel package"
	@echo "  lint-all     - Run linter on all packages"
	@echo "  clean        - Remove built binary"
	@echo "  install      - Install gum to system"
	@echo "  uninstall    - Remove gum from system"
	@echo "  run-panel   - Run panel example"
	@echo "  run-filter   - Run filter example"
	@echo "  run-choose  - Run choose example"
	@echo "  build-all   - Build for all platforms"
