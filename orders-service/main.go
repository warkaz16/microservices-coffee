package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Модель заказа
type Order struct {
	gorm.Model
	DrinkID   uint    `json:"drink_id"`
	DrinkName string  `json:"drink_name"` // Копируем из Menu Service
	Price     float64 `json:"price"`      // Копируем из Menu Service
	Quantity  int     `json:"quantity"`
	Total     float64 `json:"total"`
	Status    string  `json:"status"` // pending, completed, cancelled
}

// DTO для создания заказа
type CreateOrderRequest struct {
	DrinkID  uint `json:"drink_id"`
	Quantity int  `json:"quantity"`
}

// Структура для ответа от Menu Service
type DrinkFromMenu struct {
	ID      uint    `json:"ID"`
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	InStock bool    `json:"in_stock"`
}

var db *gorm.DB

// Адрес Menu Service
const menuServiceURL = "http://localhost:8081"

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Для локальной разработки
		dsn = "host=localhost user=adamgowz password=9555 dbname=orders_db port=5432 sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Не удалось подключиться к БД: " + err.Error())
	}

	db.AutoMigrate(&Order{})

	router := gin.Default()

	router.POST("/orders", createOrder)
	router.GET("/orders/:id", getOrder)
	router.GET("/orders", listOrders)

	router.Run(":8082")
}

// POST /orders — создать заказ
func createOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Количество должно быть больше 0"})
		return
	}

	// ====== ГЛАВНАЯ СВЯЗКА: запрос к Menu Service ======
	drink, err := getDrinkFromMenuService(req.DrinkID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Menu Service недоступен: " + err.Error()})
		return
	}

	if drink == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Напиток не найден"})
		return
	}

	if !drink.InStock {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Напитка нет в наличии"})
		return
	}
	// ====================================================

	// Создаём заказ
	order := Order{
		DrinkID:   drink.ID,
		DrinkName: drink.Name,
		Price:     drink.Price,
		Quantity:  req.Quantity,
		Total:     drink.Price * float64(req.Quantity),
		Status:    "pending",
	}

	db.Create(&order)
	c.JSON(http.StatusCreated, order)
}

// getDrinkFromMenuService — запрос к Menu Service
func getDrinkFromMenuService(drinkID uint) (*DrinkFromMenu, error) {
	url := fmt.Sprintf("%s/drinks/%d", menuServiceURL, drinkID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Напиток не найден
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	// Другая ошибка
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Menu Service вернул статус %d", resp.StatusCode)
	}

	// Читаем и парсим ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var drink DrinkFromMenu
	if err := json.Unmarshal(body, &drink); err != nil {
		return nil, err
	}

	return &drink, nil
}

// GET /orders/:id — получить заказ
func getOrder(c *gin.Context) {
	id := c.Param("id")

	var order Order
	if err := db.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заказ не найден"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GET /orders — список заказов
func listOrders(c *gin.Context) {
	var orders []Order
	db.Find(&orders)
	c.JSON(http.StatusOK, orders)
}
