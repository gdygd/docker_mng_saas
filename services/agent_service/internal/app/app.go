package app

import (
	"sync"

	"agent-service/internal/container"
	"agent-service/internal/logger"
	"agent-service/internal/server/api"
	"agent-service/internal/server/dbbatch"
	"agent-service/internal/server/gapi"
)

type Application struct {
	wg         *sync.WaitGroup
	ApiServer  *api.Server
	GApiServer *gapi.Server
	DbBatch    *dbbatch.DbBatch
}

func NewApplication(ct *container.Container, ch_terminate chan bool) *Application {
	var wg *sync.WaitGroup = &sync.WaitGroup{}

	// new dbbatch — gRPC 수신 데이터를 비동기로 DB에 기록
	batch := dbbatch.NewDbBatch(ct)

	// new http server
	apisvr, err := api.NewServer(wg, ct, ch_terminate)
	if err != nil {
		logger.Log.Error("Api server initialization fail.. %v", err)
		return nil
	}

	// new gRPC server
	gapisvr, err := gapi.NewServer(wg, ct, ch_terminate, batch)
	if err != nil {
		logger.Log.Error("gRPC server initialization fail.. %v", err)
		return nil
	}

	return &Application{
		wg:         wg,
		ApiServer:  apisvr,
		GApiServer: gapisvr,
		DbBatch:    batch,
	}
}

func (app Application) Start() {
	logger.Log.Print(3, "Start DbBatch workers..")
	app.DbBatch.Start()

	app.wg.Add(1)
	logger.Log.Print(3, "Start API server.. #1")
	go app.ApiServer.Start()

	app.wg.Add(1)
	logger.Log.Print(3, "Start gRPC server.. #1")
	go app.GApiServer.StartgPRC()
}

func (app Application) Shutdown() {
	logger.Log.Print(3, "Shutdown gRPC server#1")
	app.GApiServer.ShutdowngRPC()
	logger.Log.Print(3, "Shutdown gRPC server#2")

	// gRPC 종료 후 잔여 큐 소진
	logger.Log.Print(3, "Shutdown DbBatch#1")
	app.DbBatch.Shutdown()
	logger.Log.Print(3, "Shutdown DbBatch#2")

	logger.Log.Print(3, "Shutdown Rest server#1")
	app.ApiServer.Shutdown()
	logger.Log.Print(3, "Shutdown Rest server#2")
}
