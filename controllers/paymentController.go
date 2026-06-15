package controllers

import (
	"net/http"
	"aplikasi_admin_komunitas_badminton/db"
	"github.com/gin-gonic/gin"
	"time"
)

type GenerateBill struct {
	CommunityID    uint	   `json:"community_id" binding:"required"`
	ScheduleID     uint    `json:"schedule_id" binding:"required"`
	CourtFee       float64 `json:"court_fee"`     
	CockPriceUnit  float64 `json:"cock_price_unit"`
}

func GenerateSessionBilling(c *gin.Context) {
	var input GenerateBill
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var matches []db.Match
	if err := db.DB.Preload("CommunityPlayer").Where("schedule_id = ?", input.ScheduleID).Find(&matches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pertandingan"})
		return
	}

	if len(matches) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Belum ada pertandingan dicatat pada jadwal ini"})
		return
	}

	playerCockFeeMap := make(map[uint]float64)
	playerObjectMap := make(map[uint]db.CommunityPlayer)

	for _, m := range matches {
		totalPemain := len(m.CommunityPlayer)
		if totalPemain == 0 {
			continue
		}
		costPerPlayerInMatch := (float64(m.CocksUsed) * input.CockPriceUnit) / float64(totalPemain)

		for _, cp := range m.CommunityPlayer {
			playerCockFeeMap[cp.ID] += costPerPlayerInMatch
			playerObjectMap[cp.ID] = cp 
		}
	}

	tx := db.DB.Begin()
	var finalBills []db.Payment

	for cpID, totalBebanKok := range playerCockFeeMap {
		var existingBill db.Payment
		err := tx.Where("schedule_id = ? AND community_player_id = ?", input.ScheduleID, cpID).First(&existingBill).Error

		if err != nil {
			bill := db.Payment{
				ScheduleID:  input.ScheduleID,
				CommunityID: input.CommunityID,
				CommunityPlayerID:    cpID,
				CourtFee:    input.CourtFee, 
				CockFee:     totalBebanKok,  
				TotalPaid:   input.CourtFee + totalBebanKok,
				IsConfirmed: false,
			}
			tx.Create(&bill)
			bill.CommunityPlayer = playerObjectMap[cpID]
			finalBills = append(finalBills, bill)
		} else {
			existingBill.CourtFee = input.CourtFee
			existingBill.CockFee = totalBebanKok
			existingBill.TotalPaid = input.CourtFee + totalBebanKok
			tx.Save(&existingBill)
			existingBill.CommunityPlayer = playerObjectMap[cpID]
			finalBills = append(finalBills, existingBill)
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Tagihan iuran mabar berhasil dihitung berdasarkan pemakaian riil individu!",
		"billing_details": finalBills,
	})
}

func ConfirmPaymentOK(c *gin.Context) {
	paymentID := c.Param("id")

	var payment db.Payment
	if err := db.DB.First(&payment, paymentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tagihan tidak ditemukan"})
		return
	}

	if payment.IsConfirmed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tagihan ini sudah diverifikasi sebelumnya"})
		return
	}

	payment.IsConfirmed = true
	payment.CreatedAt = time.Now() // Tanggal pembayaran dikonfirmasi masuk kas
	db.DB.Save(&payment)

	c.JSON(http.StatusOK, gin.H{
		"message": "Pembayaran diverifikasi OK! Uang masuk kas pendapatan mabar",
		"data": payment,
	})
}