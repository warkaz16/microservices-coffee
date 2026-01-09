package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Модель напитка
type Drink struct {
	gorm.Model
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	InStock bool    `json:"in_stock"`
}

// DTO для создания напитка
type CreateDrinkRequest struct {
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	InStock bool    `json:"in_stock"`
}

var db *gorm.DB

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Для локальной разработки
		dsn = "host=localhost user=adamgowz password=9555 dbname=menu_db port=5432 sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Не удалось подключиться к БД: " + err.Error())
	}

	db.AutoMigrate(&Drink{})

	router := gin.Default()

	router.GET("/drinks", listDrinks)
	router.GET("/drinks/:id", getDrink)
	router.POST("/drinks", createDrink)

	router.Run("8081")

}

// GET /drinks — список всех напитков
func listDrinks(c *gin.Context) {
	var drinks []Drink
	db.Find(&drinks)
	c.JSON(http.StatusOK, drinks)
}

// GET /drinks/:id — один напиток
func getDrink(c *gin.Context) {
	id := c.Param("id")

	var drink Drink
	if err := db.First(&drink, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Напиток не найден"})
		return
	}

	c.JSON(http.StatusOK, drink)
}

// POST /drinks — создать напиток
func createDrink(c *gin.Context) {
	var req CreateDrinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	drink := Drink{
		Name:    req.Name,
		Price:   req.Price,
		InStock: req.InStock,
	}

	db.Create(&drink)
	c.JSON(http.StatusCreated, drink)
}
