PROTOC ?= .tools/protoc/bin/protoc.exe
PROTOC_INCLUDE ?= .tools/protoc/include
GEN_GO_DIR := gen/go
PROTO_DIRS := auth common news notifications profile requests
PROTO_FILES := $(foreach dir,$(PROTO_DIRS),$(wildcard $(dir)/v1/*.proto))

.PHONY: generate gen clean tools check-tools

generate: check-tools
	mkdir -p $(GEN_GO_DIR)
	$(PROTOC) -I . -I $(PROTOC_INCLUDE) \
		--go_out=$(GEN_GO_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_GO_DIR) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)

gen: generate

clean:
	rm -rf $(GEN_GO_DIR)

tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

check-tools:
	@$(PROTOC) --version >/dev/null
	@protoc-gen-go --version >/dev/null
	@protoc-gen-go-grpc --version >/dev/null
