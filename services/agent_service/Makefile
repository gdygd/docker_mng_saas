GENPROTO_DIR := $(shell go list -f '{{ .Dir }}' -m google.golang.org/genproto)
PROTOBUF_DIR := $(shell go list -f '{{ .Dir }}' -m google.golang.org/protobuf)

PB_OUT_DIR := ./pb/

# 테스트 실행
test:
	go test -v -cover 

start:
	cd cmd && go run main.go
	
# 빌드 실행
build:
	go build -o ./bin/agent-service ./cmd/main.go

# proto 컴파일 (grpc + gateway + annotations)
proto:
	@echo $(PROTOBUF_DIR)
	@echo $(GENPROTO_DIR)
	
	@echo "Generating Go files from proto..."
	rm -f $(PB_OUT_DIR)/*.go
	protoc \
		-I$(PROTOBUF_DIR) \
		-I$(GENPROTO_DIR) \
		--proto_path=proto \
		--go_out=$(PB_OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PB_OUT_DIR) --go-grpc_opt=paths=source_relative \
		proto/*.proto

.PHONY: build proto test
