BINARY_NAME ?= ai-ops-agent

# 默认平台
GOOS ?= linux
GOARCH ?= amd64

# 编译输出路径
OUTPUT_DIR = build/$(GOOS)_$(GOARCH)
DIST_DIR   = dist

.PHONY: build package \
        build-linux-amd64 build-linux-arm64 \
        build-windows-amd64 build-darwin-arm64 build-darwin-amd64 \
        package-linux-amd64 package-linux-arm64 \
        package-windows-amd64 package-darwin-arm64 package-darwin-amd64 \
        package-all clean

# 仅构建二进制
build:
	@echo "Building for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -o $(OUTPUT_DIR)/$(BINARY_NAME) .

# 构建并打包 tar.gz
package: build
	@mkdir -p $(DIST_DIR)
	@echo "Packaging $(BINARY_NAME)_$(GOOS)_$(GOARCH).tar.gz..."
	tar -czf $(DIST_DIR)/$(BINARY_NAME)_$(GOOS)_$(GOARCH).tar.gz -C $(OUTPUT_DIR) $(BINARY_NAME)

# 平台快捷命令（只构建）
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

# 平台快捷命令（构建 + 打 tar.gz，文件名如截图）
package-linux-amd64:
	$(MAKE) package GOOS=linux GOARCH=amd64

package-linux-arm64:
	$(MAKE) package GOOS=linux GOARCH=arm64

package-windows-amd64:
	$(MAKE) package GOOS=windows GOARCH=amd64

package-darwin-arm64:
	$(MAKE) package GOOS=darwin GOARCH=arm64

package-darwin-amd64:
	$(MAKE) package GOOS=darwin GOARCH=amd64

# 一次性打所有平台包（按需删减）
package-all: \
	package-linux-amd64 \
	package-linux-arm64 \
	package-darwin-amd64 \
	package-darwin-arm64

clean:
	rm -rf build dist