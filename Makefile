APIGW_BIN_DIR = ./services/api-gateway/bin
AUTH_BIN_DIR = ./services/auth_service/bin
AGENT_BIN_DIR = ./services/agent_service/bin

AGENT_SERVICE_PATH = ./services/agent_service
AUTH_SERVICE_PATH = ./services/auth_service
APIGW_SERVICE_PATH = ./services/api-gateway

test:
	go test -v -cover 

start:
	cd cmd && go run main.go
	
# ---------------------------------
# Build
# ---------------------------------
build-all: build-agent build-auth build-gw

build-agent:
	cd $(AGENT_SERVICE_PATH) && go build -o ./bin/ssagent-service ./cmd/main.go

build-auth:
	cd $(AUTH_SERVICE_PATH) && go build -o ./bin/ssauth-service ./cmd/main.go

build-gw:
	cd $(APIGW_SERVICE_PATH) && go build -o ./bin/ssapi-gateway ./cmd/main.go

startgw: build-gw
	cd $(APIGW_BIN_DIR)/bin && ./ssapi-gateway

startauth: build-auth
	cd $(AUTH_BIN_DIR) && ./ssauth-service


.PHONY: build