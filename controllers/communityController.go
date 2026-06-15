package controllers

import (
	"net/http"
	"aplikasi_admin_komunitas_badminton/db"
	"github.com/gin-gonic/gin"
)


func CreateCommunity(c *gin.Context) {
	var input db.Community
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data komunitas"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Komunitas berhasil didaftarkan", "data": input})
}

func GetCommunities(c *gin.Context) {
	var listCommunities []db.Community
	
	if err := db.DB.Find(&listCommunities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": listCommunities})
}

func GetCommunityByID(c *gin.Context) {
	id := c.Param("id")
	var comm db.Community
	if err := db.DB.First(&comm, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Komunitas tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": comm})
}

func UpdateCommunity(c *gin.Context) {
	var input db.Community
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingComm db.Community
	if err := db.DB.First(&existingComm, input.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Komunitas tidak ditemukan"})
		return
	}

	existingComm.Name = input.Name
	existingComm.Location = input.Location

	db.DB.Save(&existingComm)

	c.JSON(http.StatusOK, gin.H{"message": "Komunitas berhasil diperbarui", "data": existingComm})
}