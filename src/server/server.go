package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	logger "github.com/kemlee/go-rest-api-practise/core/log"

	config "github.com/kemlee/go-rest-api-practise/config"
	controller "github.com/kemlee/go-rest-api-practise/controller"
	encryptionService "github.com/kemlee/go-rest-api-practise/core/encryption"
	hashService "github.com/kemlee/go-rest-api-practise/core/hash"

	middleware "github.com/kemlee/go-rest-api-practise/core/middleware"
	userService "github.com/kemlee/go-rest-api-practise/user"
	userRepository "github.com/kemlee/go-rest-api-practise/user/repository"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	httpServer  *http.Server
	userService userService.IUserService
	logger      logger.Logger
}

func NewServer() *Server {
	logger := logger.New()

	config, _ := config.GetAPIConfig()
	db := registerDatabaseRoot(logger)
	userRepo := userRepository.GetUserRepository(db)
	hashSer := hashService.GetHashService()
	encryptionSer := encryptionService.New()
	userService := userService.GetUserService(userRepo, hashSer, encryptionSer, config)

	return &Server{
		userService: userService,
		logger:      logger,
	}
}

func (server *Server) Run() error {
	config, _ := config.GetAPIConfig()
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
		middleware.CorsMiddleware(config),
		middleware.LanguageMiddleware(config),
	)

	// Add controller
	controller.AppControllerRegister(router)
	controller.UserControllerRegister(router, server.userService)

	server.logger.Info("Server running on Port:", config.Port)
	server.httpServer = &http.Server{
		Addr:           ":" + config.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := server.httpServer.ListenAndServe(); err != nil {
			server.logger.Error("Failed to listen and serve: ", err)
		}
	}()

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, os.Interrupt, os.Interrupt)

	<-quitChan
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return server.httpServer.Shutdown(ctx)
}

func registerDatabaseRoot(logger logger.Logger) *mongo.Database {
	config, _ := config.GetAPIConfig()
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Database.Uri))
	if err != nil {
		logger.Error("Error occurred while establishing connection to mongoDB")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		logger.Error(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		logger.Error(err)
	}
	return client.Database(config.Database.Name)
}
