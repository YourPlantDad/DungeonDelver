package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Voorbeeldmodel voor een speler
type Player struct {
	gorm.Model
	Name string
	XP   int
	HP   int
	// Voeg hier extra velden toe, zoals inventory, equipment, etc.
}

func main() {
	// Stel de DSN samen op basis van de environment variables
	dsn := fmt.Sprintf(
		"host=postgres user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)

	// Maak verbinding met de database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Kon geen verbinding maken met de database:", err)
	}

	// Voer automatisch migraties uit: dit maakt de nodige tabellen aan
	if err := db.AutoMigrate(&Player{}); err != nil {
		log.Fatal("Migratie mislukt:", err)
	}

	// Initialiseer de Gin router
	r := gin.Default()

	// Voeg CORS middleware toe en sta verzoeken toe vanaf de Svelte-dev server
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type"},
	}))

	// Eenvoudige test-endpoint
	r.GET("/api/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hallo vanuit Go backend!",
		})
	})

	// Voorbeeld endpoint om een nieuwe speler aan te maken
	r.POST("/api/player", func(c *gin.Context) {
		var newPlayer Player
		if err := c.ShouldBindJSON(&newPlayer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&newPlayer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, newPlayer)
	})

	// Start de server op poort 8080
	r.Run(":8080")
}
