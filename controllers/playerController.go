package controllers

import (
	"net/http"
	"aplikasi_admin_komunitas_badminton/db"
	"github.com/gin-gonic/gin"
)

type PlayerInput struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone"`
}

func CreatePlayer(c *gin.Context) {

	var input PlayerInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPlayer := db.Player{
		Name:        input.Name,
		Phone:       input.Phone,  
	}

	if err := db.DB.Create(&newPlayer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan pemain"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pemain berhasil ditambahkan!", "data": newPlayer})
}

