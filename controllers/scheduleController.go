package controllers

import (
	"net/http"
	"time"
	"aplikasi_admin_komunitas_badminton/db"
	"github.com/gin-gonic/gin"
)

type BookCourtInput struct {
	CommunityID uint `json:"community_id" binding:"required"`
	CourtID   uint   `json:"court_id" binding:"required"`
	Type      string `json:"type" binding:"required"` 
	StartDate string `json:"start_date" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`   
	TotalCost float64 `json:"total_cost" binding:"required"`
}

func BookCourtAndSchedule(c *gin.Context) {
	var input BookCourtInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parsedDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal salah, gunakan YYYY-MM-DD"})
		return
	}

	var court db.Court
	if err := db.DB.First(&court, input.CourtID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data lapangan tidak ditemukan"})
		return
	}

	var schedules []db.Schedule
	iterations := 1
	
	if input.Type == "bulanan-4" {
		iterations = 4 
	}else if input.Type == "bulanan-5"{
		iterations = 5
	}

	tx := db.DB.Begin()

	currentDate := parsedDate
	for i := 0; i < iterations; i++ {
		newSchedule := db.Schedule{
			CourtID:     input.CourtID,
			Date:        currentDate,
			StartTime:   input.StartTime,
			EndTime:     input.EndTime,
			Status:      "Booked",
		}
		
		if err := tx.Create(&newSchedule).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat jadwal pekan ke-"})
			return
		}
		schedules = append(schedules, newSchedule)

		currentDate = currentDate.AddDate(0, 0, 7)
	}

	notes := "Bayar sewa " + court.Name + " (" + input.Type + ") mulai tanggal " + input.StartDate
	receipt := db.Receipt{
		CommunityID: input.CommunityID,
		Category:    "Sewa Lapangan",
		Amount:      input.TotalCost,
		Notes:       notes,
		Date:        time.Now(),
	}

	if err := tx.Create(&receipt).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat pengeluaran lapangan"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Jadwal berhasil terbentuk dan pengeluaran dicatat",
		"schedules": schedules,
		"receipt": receipt,
	})
}