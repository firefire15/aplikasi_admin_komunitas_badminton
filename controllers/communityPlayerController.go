package controllers

import (
	"net/http"
	"aplikasi_admin_komunitas_badminton/db"
	"github.com/gin-gonic/gin"
)

type CommunityPlayerInput struct{
	CommunityID uint `json:"community_id" binding:"required"`
	PlayerID uint `json:"player_id" binding:"required"`
	Status		string `json:"status" binding:"required"`
}

func AddCommunityPlayers(c *gin.Context) {
	var input CommunityPlayerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newCommunityPlayer := db.CommunityPlayer{
		CommunityID : input.CommunityID,
		PlayerID: input.PlayerID,
		Status: input.Status,
	}

	if err := db.DB.Create(&newCommunityPlayer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data pemain komunitas"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pemain komunitas berhasil didaftarkan", "data": input})
}

func GetMyCommunityPlayers(c *gin.Context) {
	communityID := c.Param("id")

	var communityPlayer []db.CommunityPlayer
	db.DB.Preload("Community").Preload("Player").Where("community_id = ?", communityID).Find(&communityPlayer)

	c.JSON(http.StatusOK, gin.H{"data": communityPlayer})
}