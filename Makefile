BINARY_NAME ?= ai-ops-agent

# 模块路径：改成你 go.mod 里的 module 名
MODULE_PATH ?= ai-ops-agent

# 默认平台
GOOS  ?= linux
GOARCH ?= amd64

OUTPUT_DIR = build/$(GOOS)_$(GOARCH)
DIST_DIR   = dist

# 版本信息（从 git 自动取，取不到就给个兜底）
VERSION  := $(shell git describe --tags --abbrev=0 2>/dev/null || echo dev)
COMMIT   := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
BUILDDATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

LDFLAGS := -X '$(MODULE_PATH)/internal/version.Version=$(VERSION)' \
           -X '$(MODULE_PATH)/internal/version.Commit=$(COMMIT)' \
           -X '$(MODULE_PATH)/internal/version.BuildDate=$(BUILDDATE)'

.PHONY: build package \
        build-linux-amd64 build-linux-arm64 \
        build-darwin-amd64 build-darwin-arm64 \
        package-linux-amd64 package-linux-arm64 \
        package-darwin-amd64 package-darwin-arm64 \
        package-all clean

# build only
build:
	@mkdir -p $(OUTPUT_DIR)
	@echo "Building $(BINARY_NAME) $(VERSION) for $(GOOS)/$(GOARCH)..."
	@GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 \
		go build -ldflags "$(LDFLAGS)" -o $(OUTPUT_DIR)/$(BINARY_NAME) .

# build + tar.gz (完全静默)
package: build
	@mkdir -p $(DIST_DIR)
	@COPYFILE_DISABLE=1 tar zcf $(DIST_DIR)/$(BINARY_NAME)_$(GOOS)_$(GOARCH).tar.gz -C $(OUTPUT_DIR) $(BINARY_NAME)

# shortcuts
package-linux-amd64:
	$(MAKE) package GOOS=linux GOARCH=amd64

package-linux-arm64:
	$(MAKE) package GOOS=linux GOARCH=arm64

package-darwin-amd64:
	$(MAKE) package GOOS=darwin GOARCH=amd64

package-darwin-arm64:
	$(MAKE) package GOOS=darwin GOARCH=arm64

package-all: \
	package-linux-amd64 \
	package-linux-arm64 \
	package-darwin-amd64 \
	package-darwin-arm64

clean:
	rm -rf build dist