package main

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Маршруты к Menu Service
	router.GET("/api/drinks", getDrinks)
	router.GET("/api/drinks/:id", getDrinkByID)
	router.POST("/api/drinks", createDrink)

	// Маршруты к Orders Service
	router.GET("/api/orders", getOrders)
	router.GET("/api/orders/:id", getOrderByID)
	router.POST("/api/orders", createOrder)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Для локальной разработки
	}
	router.Run(":" + port)
}

// ========== Menu Service ==========

func getDrinks(c *gin.Context) {
	resp, err := http.Get("http://localhost:8081/drinks")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Menu Service недоступен"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

func getDrinkByID(c *gin.Context) {
	id := c.Param("id")

	resp, err := http.Get("http://localhost:8081/drinks/" + id)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Menu Service недоступен"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

func createDrink(c *gin.Context) {
	// Читаем тело запроса от клиента
	reqBody, _ := io.ReadAll(c.Request.Body)

	// Пересылаем в Menu Service
	resp, err := http.Post(
		"http://localhost:8081/drinks",
		"application/json",
		bytes.NewReader(reqBody),
	)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Menu Service недоступен"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

// ========== Orders Service ==========

func getOrders(c *gin.Context) {
	resp, err := http.Get("http://localhost:8082/orders")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Orders Service недоступен"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

func getOrderByID(c *gin.Context) {
	id := c.Param("id")

	resp, err := http.Get("http://localhost:8082/orders/" + id)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Orders Service недоступен"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

func createOrder(c *gin.Context) {
	reqBody, _ := io.ReadAll(c.Request.Body)

	resp, err := http.Post(
		"http://localhost:8082/orders",
		"application/json",
		bytes.NewReader(reqBody),
	)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Orders Service недоступен"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}
