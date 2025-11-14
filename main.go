package main

import (
	"fmt"
	"hackathon/database"
	"hackathon/pkg"
	"hackathon/user"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Container struct {
	DB          *gorm.DB
	UserHandler *user.Handler
	UserService *user.Service
	UserRepo    *user.Repository
	JWTService  *pkg.JWTService
}

func NewContainer() (*Container, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("godotenv.Load: %w", err)
	}

	db, err := database.Connect()
	if err != nil {
		return nil, fmt.Errorf("database.Connect: %w", err)
	}

	accessSecret := getEnv("ACCESS_SECRET")
	refreshSecret := getEnv("REFRESH_SECRET")

	jwtService := pkg.NewJWTService(
		[]byte(accessSecret),
		[]byte(refreshSecret),
		15*time.Minute,
		7*24*time.Hour,
	)

	userRepo := user.NewRepository(db)

	userService := user.NewUserService(userRepo, jwtService)

	userHandler := user.NewHandler(userService)

	return &Container{
		DB:          db,
		UserHandler: userHandler,
		UserService: userService,
		UserRepo:    userRepo,
		JWTService:  jwtService,
	}, nil
}

func (c *Container) RunMigrations() error {
	models := []any{
		&user.Model{},
	}

	if err := c.DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("c.DB.AutoMigrate: %w", err)
	}

	log.Println("✅ Database migrations completed")
	return nil
}

func (c *Container) SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	auth := router.Group("/auth")
	{
		auth.POST("/signup", c.UserHandler.SignUp)
		auth.POST("/signin", c.UserHandler.SignIn)
		// refresh logout
	}

	// api := router.Group("/api")

	return router
}

func getEnv(key string) string {
	return os.Getenv(key)
}

func main() {
	container, err := NewContainer()
	if err != nil {
		log.Fatalf("❌ Failed to initialize container: %v", err)
	}

	if err := container.RunMigrations(); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	router := container.SetupRouter()

	if err := router.Run(":9999"); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
