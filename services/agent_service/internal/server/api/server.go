package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"agent-service/internal/config"
	"agent-service/internal/container"
	"agent-service/internal/db"
	"agent-service/internal/logger"
	"agent-service/internal/memory"
	"agent-service/internal/service"

	apiserv "agent-service/internal/service/api"

	"github.com/gdygd/goglib/token"

	"github.com/gin-gonic/gin"
)

const (
	R_TIME_OUT = 5 * time.Second
	W_TIME_OUT = 5 * time.Second
)

// Server serves HTTP requests for our banking service.
type Server struct {
	wg           *sync.WaitGroup
	srv          *http.Server
	config       *config.Config
	tokenMaker   token.Maker
	router       *gin.Engine
	service      service.ServiceInterface
	dbHnd        db.DbHandler
	objdb        *memory.RedisDb
	ch_terminate chan bool
}

func NewServer(wg *sync.WaitGroup, ct *container.Container, ch_terminate chan bool) (*Server, error) {
	// init service
	apiservice := apiserv.NewApiService(ct.DbHnd, ct.ObjDb)
	tokenMaker, err := token.NewJWTMaker(ct.Config.TokenSecretKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker:%w", err)
	}

	server := &Server{
		wg:           wg,
		config:       ct.Config,
		tokenMaker:   tokenMaker,
		service:      apiservice,
		dbHnd:        ct.DbHnd,
		objdb:        ct.ObjDb,
		ch_terminate: ch_terminate,
	}

	server.setupRouter()

	server.srv = &http.Server{}
	server.srv.Addr = ct.Config.HTTPServerAddress
	server.srv.Handler = server.router.Handler()
	server.srv.ReadTimeout = R_TIME_OUT
	server.srv.WriteTimeout = W_TIME_OUT

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	// router := gin.New()
	// addresses := strings.Split(server.config.AllowOrigins, ",")
	// router.Use(corsMiddleware(addresses))
	// router.Use(authMiddleware(server.tokenMaker))

	// gin.SetMode(gin.DebugMode)
	// fmt.Printf("%v, \n", server.config.AllowOrigins)

	router.GET("/heartbeat", server.heartbeat)
	router.GET("/terminate", server.terminate)

	router.GET("/test", server.testapi)

	router.GET("/ps2/:agentid/:hostid", server.dockerPs2) // id: agent id, host : host id
	router.GET("/inspect2/:agentid/:host/:id", server.containerInspect2)
	router.GET("/stat3/:agentid/:host", server.statContainer3)

	router.GET("/events", gin.WrapF(handleSSE()))

	server.router = router
}

func (server *Server) Start() error {
	logger.Log.Print(2, "Gin server start.")

	// sse
	go ProcessEventMsg()

	if err := server.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Error("listen error. %v", err)
		return err
	}

	return nil
}

func (server *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer server.wg.Done()
	if err := server.srv.Shutdown(ctx); err != nil {
		logger.Log.Error("Server Shutdown:", err)
		return err
	}
	return nil
}

func SetCookie(c *gin.Context, name string, value string, expTm time.Time, maxAge int) {
	// 운영환경 여부 (예: ENV=prod)
	isProd := os.Getenv("ENV") == "prod"

	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expTm,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   isProd, // 운영환경이면 HTTPS만 허용
		SameSite: http.SameSiteLaxMode,
	}

	// Domain은 환경변수로 설정 (IP 환경이면 비워둠)
	domain := os.Getenv("COOKIE_DOMAIN")
	if domain != "" {
		cookie.Domain = domain
	}

	http.SetCookie(c.Writer, cookie)
}

/*
exp := time.Now().Add(7 * 24 * time.Hour)

SetCookie(
	c,
	"refreshToken",
	refreshToken,
	exp,
	int((7 * 24 * time.Hour).Seconds()),
)
*/
