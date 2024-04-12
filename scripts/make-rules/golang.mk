# ==============================================================================
# Makefile helper functions for golang
#

GO := go

GO_BUILD_FLAGS += -ldflags "$(GO_LDFLAGS)"

ifeq ($(GOOS),windows)
	GO_OUT_EXT := .exe
endif

ifeq ($(ROOT_PACKAGE),)
	$(error the variable ROOT_PACKAGE must be set prior to including golang.mk)
endif

GOPATH := $(shell go env GOPATH)
ifeq ($(origin GOBIN), undefined)
	GOBIN := $(GOPATH)/bin
endif

COMMANDS ?= $(filter-out %.md, $(wildcard $(ROOT_DIR)/cmd/*))
BINS ?= $(foreach cmd,$(COMMANDS),$(notdir $(cmd)))

ifeq ($(COMMANDS),)
  $(error Could not determine COMMANDS, set ROOT_DIR or run in source dir)
endif
ifeq ($(BINS),)
  $(error Could not determine BINS, set ROOT_DIR or run in source dir)
endif

.PHONY: go.build.verify
go.build.verify:
	@if ! which go &>/dev/null; then echo "Cannot found go compile tool. Please install go tool first."; exit 1; fi

.PHONY: go.build.%
go.build.%:
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> Building binary $(COMMAND) $(VERSION) for $(OS) $(ARCH)"
	@mkdir -p $(OUTPUT_DIR)/bin
	@CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) $(GO) build $(GO_BUILD_FLAGS) -o $(OUTPUT_DIR)/bin/$(COMMAND)$(GO_OUT_EXT) $(ROOT_PACKAGE)/cmd/$(COMMAND)

.PHONY: go.build
go.build: go.build.verify $(addprefix go.build., $(addprefix $(PLATFORM)., $(BINS)))

.PHONY: go.format
go.format: tools.verify.goimports ## 格式化 Go 源码.
	@$(FIND) -type f -name '*.go' | $(XARGS) gofmt -s -w
	@$(FIND) -type f -name '*.go' | $(XARGS) goimports -w -local $(ROOT_PACKAGE)
	@$(GO) mod edit -fmt

.PHONY: go.tidy
go.tidy: ## 自动添加/移除依赖包.
	@$(GO) mod tidy

.PHONY: go.test
go.test: ## 执行单元测试.
	@echo "===========> Run unit test"
	@mkdir -p $(OUTPUT_DIR)
	@set -o pipefail;$(GO) test -race -cover -coverprofile=$(OUTPUT_DIR)/coverage.out -timeout=10m -shuffle=on -short -v `go list ./...`
	# @sed -i '/mock_.*.go/d' $(OUTPUT_DIR)/coverage.out # 从 coverage 中删除mock_.*.go 文件
	# @sed -i '/internal\/apiserver\/store\/.*.go/d' $(OUTPUT_DIR)/coverage.out # store/ 下的 Go 代码不参与覆盖率计算

.PHONY: go.cover
go.cover: go.test ## 执行单元测试，并校验覆盖率阈值.
	# @$(GO) tool cover -func=$(OUTPUT_DIR)/coverage.out | awk -v target=$(COVERAGE) -f $(ROOT_DIR)/scripts/coverage.awk
	@$(GO) tool cover -func=$(OUTPUT_DIR)/coverage.out

.PHONY: go.lint
go.lint: tools.verify.golangci-lint
	@echo "===========> Run golangci to lint source codes"
	@golangci-lint run -c $(ROOT_DIR)/.golangci.yaml $(ROOT_DIR)/...
