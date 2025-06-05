# ==============================================================================
# 定义全局 Makefile 变量方便后面引用

COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
# 项目根目录
PROJ_ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/ && pwd -P))
# 构建产物、临时文件存放目录
OUTPUT_DIR := $(PROJ_ROOT_DIR)/_output

# 将 Makefile 中的 Shell 切换为 bash
SHELL := /bin/bash


# ==== 配置 ====
DEP_DIR = deps

PROTOC_VERSION = 31.1

# 自动判断平台架构
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

ifeq ($(UNAME_S),Darwin)
	ifeq ($(UNAME_M),arm64)
		PROTOC_PLATFORM = osx-aarch_64
	else
		PROTOC_PLATFORM = osx-x86_64
	endif
else ifeq ($(UNAME_S),Linux)
	ifeq ($(UNAME_M),aarch64)
		PROTOC_PLATFORM = linux-aarch_64
	else
		PROTOC_PLATFORM = linux-x86_64
	endif
endif

PROTOC_ZIP = protoc-$(PROTOC_VERSION)-$(PROTOC_PLATFORM).zip
PROTOC_URL = https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/$(PROTOC_ZIP)

# 设置本地 bin 目录优先
export PATH := $(abspath $(BIN_DIR)):$(PATH)

# ==== 安装 protoc 到 dep/ ====
install-protoc:
	@mkdir -p $(DEP_DIR)
	@echo "Detected platform: $(PROTOC_PLATFORM)"
	@echo "Downloading protoc..."
	curl -LO $(PROTOC_URL)
	unzip -j $(PROTOC_ZIP) bin/protoc -d $(DEP_DIR)
	rm -rf $(PROTOC_ZIP)
	@echo "✅ protoc installed to $(DEP_DIR)/protoc"
	@$(DEP_DIR)/protoc --version

# ==== 安装插件到 dep ====
install-plugins:
	GOBIN=$(abspath $(DEP_DIR)) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(abspath $(DEP_DIR)) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(abspath $(DEP_DIR)) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(abspath $(DEP_DIR)) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# ==== 检查插件 ====
check-plugins:
	@$(DEP_DIR)/protoc --version || (echo "❌ protoc not found"; exit 1)
	@command -v $(DEP_DIR)/protoc-gen-go >/dev/null || echo "❌ protoc-gen-go not found"
	@command -v $(DEP_DIR)/protoc-gen-go-grpc >/dev/null || echo "❌ protoc-gen-go-grpc not found"
	@command -v $(DEP_DIR)/protoc-gen-grpc-gateway >/dev/null || echo "❌ protoc-gen-grpc-gateway not found"
	@command -v $(DEP_DIR)/protoc-gen-openapiv2 >/dev/null || echo "❌ protoc-gen-openapiv2 not found"
	@echo "✅ All plugins are installed in $(DEP_DIR)"

# ==== Buf 代码生成 ====
proto: buf.dep-update buf.generate

buf.generate:
	@PATH=$(abspath $(DEP_DIR)):$$PATH buf generate -v
buf.dep-update:
	@PATH=$(abspath $(DEP_DIR)):$$PATH buf dep update -v
buf.lint:
	@PATH=$(abspath $(DEP_DIR)):$$PATH buf lint -v

# ==== 启动服务 ====


# ==== 一键初始化 ====
init: install-protoc install-plugins check-plugins proto

.PHONY: init proto run install-protoc install-plugins check-plugins


# ==============================================================================
# 定义版本相关变量

## 指定应用使用的 version 包，会通过 `-ldflags -X` 向该包中指定的变量注入值
VERSION_PACKAGE=minifast/pkg/common/version
## 定义 VERSION 语义化版本号
ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --tags --always --match='v*')
endif

## 检查代码仓库是否是 dirty（默认dirty）
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
    GIT_TREE_STATE="clean"
endif
GIT_COMMIT:=$(shell git rev-parse HEAD)

GO_LDFLAGS += \
    -X $(VERSION_PACKAGE).gitVersion=$(VERSION) \
    -X $(VERSION_PACKAGE).gitCommit=$(GIT_COMMIT) \
    -X $(VERSION_PACKAGE).gitTreeState=$(GIT_TREE_STATE) \
    -X $(VERSION_PACKAGE).buildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

.PHONY: build
build: tidy # 编译源码，依赖 tidy 目标自动添加/移除依赖包.
	@go build -v -ldflags "$(GO_LDFLAGS)" -o $(OUTPUT_DIR)/mf-userserver $(PROJ_ROOT_DIR)/cmd/user/user.go

.PHONY: tidy
tidy: # 自动添加/移除依赖包.
	@go mod tidy
