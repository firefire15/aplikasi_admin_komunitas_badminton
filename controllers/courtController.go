package controllers

import (
	"net/http"
	"aplikasi_admin_komunitas_badminton/db"
	"github.com/gin-gonic/gin"
)

type CourtInput struct{
	CommunityID uint `json:"community_id" binding:"required"`
	Name string `json:"name" binding:"required"`
	PricePerDay	float64 `json:"price_per_day" binding:"required"`
}

func AddCourt(c *gin.Context) {
	var input CourtInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newCourt := db.Court{
		CommunityID : input.CommunityID,
		Name: input.Name,
		PricePerDay: input.PricePerDay,
	}

	if err := db.DB.Create(&newCourt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data lapangan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Data Lapangan berhasil didaftarkan", "data": input})
}

func GetCourt(c *gin.Context) {
	communityID := c.Param("id")

	var court []db.Court
	db.DB.Preload("Community").Where("community_id = ?", communityID).Find(&court)

	c.JSON(http.StatusOK, gin.H{"data": court})
}