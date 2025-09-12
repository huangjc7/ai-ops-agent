# Makefile

BINARY_NAME = ai

# 默认平台
GOOS ?= linux
GOARCH ?= amd64

# 编译输出路径
OUTPUT_DIR = build/$(GOOS)_$(GOARCH)

# 构建命令
build:
	@echo "Building for $(GOOS)/$(GOARCH)..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -o $(OUTPUT_DIR)/$(BINARY_NAME) .

# 快捷命令
build-linux-amd64:
	$(MAKE) build GOOS=linux GOARCH=amd64

build-linux-arm64:
	$(MAKE) build GOOS=linux GOARCH=arm64

build-windows-amd64:
	$(MAKE) build GOOS=windows GOARCH=amd64

build-darwin-arm64:
	$(MAKE) build GOOS=darwin GOARCH=arm64

build-darwin-amd64:
	$(MAKE) build GOOS=darwin GOARCH=amd64

clean:
	rm -rf build