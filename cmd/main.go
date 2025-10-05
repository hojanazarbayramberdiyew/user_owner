// @title           Test API
// @version         1.0
// @description     Test Swagger
// @host            localhost:8080
// @BasePath        /api/v1
package main

import (
	"context"
	"fmt"
	"log"
	"user_owner/internal/config"
	"user_owner/internal/handler"
	"user_owner/internal/middleware"
	"user_owner/internal/repository"
	"user_owner/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "user_owner/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	db := connectDB()
	defer db.Close()
	cfg := config.GetConfig()

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cfg)
	userHandler := handler.NewUserHandler(userService)

	api := r.Group("/api/v1")
	{
		api.POST("/users", userHandler.CreateUser)
		api.GET("/users", userHandler.GetAllUsers)
		api.POST("/login", userHandler.Login)
	}

	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	protected := r.Group("/api/v1/protected")
	protected.Use(middleware.AuthMiddleware(cfg, userRepo))
	{
		protected.POST("/orders", orderHandler.CreateOrder)
	}

	port := fmt.Sprintf(":%s", cfg.Listen.Port)
	err := r.Run(port)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}

}

func connectDB() *pgxpool.Pool {

	cfg := config.GetConfig()

	ctx := context.Background()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Storage.Username, cfg.Storage.Password, cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.DbName)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("PostgreSQL bilen hakykatdan baglanyşyk başa barmady: %v\n", err)
	}

	fmt.Println("Successfully connected!")
	return pool

}
