package main

import (
	"fmt"
	"hackathon/database"
	"hackathon/models"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Container struct {
	DB *gorm.DB
}

func NewContainer() (*Container, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("godotenv.Load: %w", err)
	}

	db, err := database.Connect()
	if err != nil {
		return nil, fmt.Errorf("database.Connect: %w", err)
	}

	return &Container{DB: db}, nil
}

func (c *Container) CreateTables() error {
	if err := c.DB.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("c.DB.AutoMigrate: %w", err)
	}
	return nil
}

func main() {
	r := gin.Default()

	container, err := NewContainer()
	if err != nil {
		log.Fatalf("NewContainer: %v", err)
	}

	if err := container.CreateTables(); err != nil {
		log.Fatalf("container.CreateTables: %v", err)
	}

	if err := r.Run(":9999"); err != nil {
		log.Fatalf("r.Run: %v", err)
	}
}
