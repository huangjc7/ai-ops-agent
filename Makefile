# 默认目标
.PHONY: all build-arm64 build-amd64 clean

# 输出目录
OUTPUT_DIR = bin

# Go 文件路径
MAIN_FILE = ai-ops-agent.go

# 目标二进制文件名
ARM_TARGET = $(OUTPUT_DIR)/ai-ops-agent-arm64
X86_TARGET = $(OUTPUT_DIR)/ai-ops-agent-amd64

# all 目标，依赖于 ARM 和 x86 的编译目标
all: build-arm64 build-amd64

# 编译 ARM 架构
build-arm64:
	mkdir -p $(OUTPUT_DIR)
	GOOS=linux GOARCH=arm64 go build -o $(ARM_TARGET) $(MAIN_FILE)

# 编译 x86 架构
build-amd64:
	mkdir -p $(OUTPUT_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(X86_TARGET) $(MAIN_FILE)

# 清理生成的文件
clean:
	rm -rf $(OUTPUT_DIR)