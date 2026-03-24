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
# 	cd $(AGENT_SERVICE_PATH) && go build -o ./bin/agent-service ./cmd/main.go
	cd $(AGENT_SERVICE_PATH) && make build

build-auth:
# 	cd $(AUTH_SERVICE_PATH) && go build -o ./bin/auth-service ./cmd/main.go
	cd $(AUTH_SERVICE_PATH) && make build

build-gw:
# 	cd $(APIGW_SERVICE_PATH) && go build -o ./bin/api-gateway ./cmd/main.go
	cd $(APIGW_SERVICE_PATH) && make build

startgw: build-gw
	cd $(APIGW_BIN_DIR) && ./api-gateway

startauth: build-auth
	cd $(AUTH_BIN_DIR) && ./auth-service

startagent: build-agent
	cd $(AGENT_BIN_DIR) && ./agent-service

allstart: 
	cd $(APIGW_SERVICE_PATH)/bin && ./api-gateway &
	cd $(AUTH_BIN_DIR) && ./auth-service &
	cd $(AGENT_BIN_DIR) && ./agent-service &

# 포트 기반 종료 (더 확실함)
allstop-port:
	@echo "Stopping all services by port..."	
	@fuser -k 19091/tcp 2>/dev/null || true
	@fuser -k 19190/tcp 2>/dev/null || true
	@fuser -k 19081/tcp 2>/dev/null || true
	@fuser -k 19082/tcp 2>/dev/null || true
	@fuser -k 19083/tcp 2>/dev/null || true
	@fuser -k 19192/tcp 2>/dev/null || true
	@echo "All services stopped"

proto-agent:
	cd $(AGENT_SERVICE_PATH) && make proto

.PHONY: build