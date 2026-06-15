package controllers

import (
	"net/http"
	"time"
	"aplikasi_admin_komunitas_badminton/db"
	"github.com/gin-gonic/gin"
)


type BuyShuttlecockInput struct {
	CommunityID uint `json:"community_id" binding:"required"`
	Brand     string  `json:"brand" binding:"required"`    
	QtySlop   int     `json:"qty_slop" binding:"required"`  
	PricePerSlop float64 `json:"price_per_slop" binding:"required"`
}

type ReturnCockInput struct {
	ShuttlecockID uint `json:"shuttlecock_id" binding:"required"`
	Qty       int  `json:"qty_ball" binding:"required"`
}

func BuyShuttlecock(c *gin.Context) {

	var input BuyShuttlecockInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	totalBalls := input.QtySlop * 12 
	totalCost := float64(input.QtySlop) * input.PricePerSlop

	tx := db.DB.Begin()

	var cock db.Shuttlecock
	err := tx.Where("community_id = ? AND brand = ?", input.CommunityID, input.Brand).First(&cock).Error

	if err == nil {
		cock.Stock += totalBalls
		tx.Save(&cock)
	} else {
		cock = db.Shuttlecock{
			CommunityID: input.CommunityID,
			Brand:       input.Brand,
			Stock:   totalBalls,
		}
		tx.Create(&cock)
	}

	notes := "Pembelian Shuttlecock " + input.Brand + " sebanyak " + string(rune(input.QtySlop)) + " slop."
	receipt := db.Receipt{
		CommunityID: input.CommunityID,
		Category:    "Beli Shuttlecock",
		Amount:      totalCost,
		Notes:       notes,
		Date:        time.Now(),
	}

	if err := tx.Create(&receipt).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat pengeluaran kok"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Stok kok berhasil ditambahkan & pengeluaran dicatat",
		"current_stock_balls": cock.Stock,
		"total_spent": totalCost,
	})
}

func GetShuttlecockStock(c *gin.Context) {
	communityID := c.Param("community_id")

	var stocks []db.Shuttlecock
	db.DB.Where("community_id = ?", communityID).Find(&stocks)

	c.JSON(http.StatusOK, gin.H{"shuttlecock_inventory": stocks})
}

func ReturnShuttlecock(c *gin.Context) {
	var input ReturnCockInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := db.DB.Begin()

	var cock db.Shuttlecock
	if err := tx.First(&cock, input.ShuttlecockID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Data kok tidak ditemukan"})
		return
	}

	if cock.Stock < input.Qty {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stok tidak mencukupi untuk diretur"})
		return
	}

	cock.Stock -= input.Qty
	cock.Returned += input.Qty
	tx.Save(&cock)

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mencatat retur kok",
		"brand": cock.Brand,
		"current_stock": cock.Stock,
		"total_returned": cock.Returned,
	})
}