package main

import (
	"encoding/json"
	"fmt"
	"os"
	"pizzaria/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

var pizzas []models.Pizza

func main() {
	loadPizzas()
	router := gin.Default()
	router.GET("/pizzas", getPizzas)
	router.POST("/pizzas", postPizzas)
	router.GET("/pizzas/:id", getPizzasById)
	router.Run() // listening on 0.0.0.0:8080
}

func getPizzas(c *gin.Context) {
	c.JSON(200, gin.H{
		"pizzas": pizzas,
	})
}

func postPizzas(c *gin.Context) {
	var newPizza models.Pizza
	if err := c.ShouldBindJSON(&newPizza); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	newPizza.Id = len(pizzas) + 1
	pizzas = append(pizzas, newPizza)
	savePizza()
	c.JSON(201, newPizza)
}

func getPizzasById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam) // necessario a conversao para int para comparar com os da lista
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}
	for _, pizza := range pizzas {
		if pizza.Id == id {
			c.JSON(200, pizza)
			return
		}
	}
	c.JSON(404, gin.H{"error": "Pizza not found"})
}

func loadPizzas() {
	file, err := os.Open("dados/pizzas.json")
	if err != nil {
		fmt.Println("Error file:", err)
		return
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&pizzas); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
}

func savePizza() {
	file, err := os.Create("dados/pizzas.json")
	if err != nil {
		fmt.Println("Error file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(pizzas); err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
}
