package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
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
}

func NewServer() *Server {
	db := registerDatabaseRoot()
	userRepo := userRepository.GetUserRepository(db)
	hashSer := hashService.GetHashService()
	encryptionSer := encryptionService.New()
	userService := userService.GetUserService(userRepo, hashSer, encryptionSer)

	return &Server{
		userService: userService,
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

	log.Printf("Server running on %v \n", config.Port)
	server.httpServer = &http.Server{
		Addr:           ":" + config.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := server.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, os.Interrupt, os.Interrupt)

	<-quitChan
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return server.httpServer.Shutdown(ctx)
}

func registerDatabaseRoot() *mongo.Database {
	config, _ := config.GetAPIConfig()
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Database.Uri))
	if err != nil {
		log.Fatalf("Error occurred while establishing connection to mongoDB")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	return client.Database(config.Database.Name)
}
