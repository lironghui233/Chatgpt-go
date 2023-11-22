package main

import (
	"chatgpt-web/pkg/cmd"
	"chatgpt-web/pkg/config"
	"chatgpt-web/pkg/log"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"chatgpt-web/pkg/controllers"
	"chatgpt-web/pkg/middlewares"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

type ChatGPTWebServer struct {
	config *config.Config
	log    log.ILogger
}

func NewChatGPTWebServer(config *config.Config, log log.ILogger) *ChatGPTWebServer {
	return &ChatGPTWebServer{
		config: config,
		log:    log,
	}
}

func (r *ChatGPTWebServer) Run(ctx context.Context) error {
	gin.SetMode(gin.ReleaseMode)
	go r.httpServer(ctx)
	return nil
}

func (r *ChatGPTWebServer) httpServer(ctx context.Context) {
	chatService, err := controllers.NewChatService(r.config, r.log)
	if err != nil {
		klog.Fatal(err)
	}

	addr := fmt.Sprintf("%s:%d", r.config.Http.Host, r.config.Http.Port)
	r.log.InfoF("ChatGPT Web Server on: %s", addr)
	server := &http.Server{
		Addr: addr,
	}
	entry := gin.Default()
	entry.Use(middlewares.Cors())
	chat := entry.Group("/api")
	if len(r.config.Http.BasicAuthUser) > 0 {
		accounts := gin.Accounts{}
		users := strings.Split(r.config.Http.BasicAuthUser, ",")
		passwords := strings.Split(r.config.Http.BasicAuthPassword, ",")
		if len(users) != len(passwords) {
			panic("basic auth setting error")
		}
		for i := 0; i < len(users); i++ {
			accounts[users[i]] = passwords[i]
		}
		chat.POST("/chat-process", gin.BasicAuth(accounts), middlewares.RateLimitMiddleware(1, 2), chatService.ChatProcess)
	} else {
		chat.POST("/chat-process", middlewares.RateLimitMiddleware(1, 2), chatService.ChatProcess)
	}
	chat.POST("/config", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "Success",
			"data": map[string]string{
				"apiModel":   "ChatGPTAPI",
				"socksProxy": "",
			},
		})
	})
	chat.POST("/session", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":  "Success",
			"message": "",
			"data": gin.H{
				"auth": false,
			},
		})
	})
	chat.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{})
	})

	entry.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{})
	})

	server.Handler = entry
	go func(ctx context.Context) {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.InfoF("Server shutdown with error %v", err)
		}
	}(ctx)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.FatalF("Server listen and serve error %v", err)
	}
}

func main() {
	loadDependOn()
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Kill, os.Interrupt)
	defer stop()

	cnf := config.GetConf()
	logger := log.NewLogger()
	logger.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	logger.SetLevel(cnf.Log.Level)
	logger.SetPrintCaller(true)

	app := NewChatGPTWebServer(cnf, logger)
	err := app.Run(ctx)
	log.Error(err)
	<-ctx.Done()
}

func loadDependOn() {
	config.InitConf(cmd.Args.Config)
	cnf := config.GetConf()
	log.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	log.SetLevel(cnf.Log.Level)
	log.SetPrintCaller(true)
}
