package controllers

import (
	"net/http"
	"aplikasi_admin_komunitas_badminton/db"
	"github.com/gin-gonic/gin"
)

type RecordMatchInput struct {
	ScheduleID     uint   `json:"schedule_id" binding:"required"`
	FormatScore    string `json:"format_score"` 
	ShuttlecockID  uint   `json:"shuttlecock_id" binding:"required"`
	CocksUsed      int    `json:"cocks_used" binding:"required"`
	DurationMinute int    `json:"duration_minute"`
	ScoreSet1      string `json:"score_set1" binding:"required"`
	ScoreSet2      string `json:"score_set2"`
	ScoreSet3      string `json:"score_set3"`
	PlayerIDs      []uint `json:"player_ids" binding:"required"`
}

type UpdateMatchInput struct {
	FormatScore    string `json:"format_score"`
	ShuttlecockID  uint   `json:"shuttlecock_id" binding:"required"`
	CocksUsed      int    `json:"cocks_used" binding:"required"`
	DurationMinute int    `json:"duration_minute"`
	ScoreSet1      string `json:"score_set1" binding:"required"`
	ScoreSet2      string `json:"score_set2"`
	ScoreSet3      string `json:"score_set3"`
	PlayerIDs      []uint `json:"player_ids" binding:"required"`
}

func RecordMatch(c *gin.Context) {
	var input RecordMatchInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx := db.DB.Begin()
	var cock db.Shuttlecock
	if err := tx.First(&cock, input.ShuttlecockID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Stok kok tidak ditemukan"})
		return
	}
	if cock.Stock < input.CocksUsed {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stok kok tidak mencukupi"})
		return
	}

	cock.Stock -= input.CocksUsed
	tx.Save(&cock)

	var Communityplayers []db.CommunityPlayer
	tx.Where("id IN ?", input.PlayerIDs).Find(&Communityplayers)
	totalPlayersInMatch := len(Communityplayers)

	if totalPlayersInMatch != 2 && totalPlayersInMatch != 4 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jumlah pemain harus 2 (Single) atau 4 (Double)"})
		return
	}

	newMatch := db.Match{
		ScheduleID:     input.ScheduleID,
		FormatScore:    input.FormatScore,
		CocksUsed:      input.CocksUsed,
		ShuttlecockID:  input.ShuttlecockID,
		DurationMinute: input.DurationMinute,
		ScoreSet1:      input.ScoreSet1,
		ScoreSet2:      input.ScoreSet2,
		ScoreSet3:      input.ScoreSet3,
		CommunityPlayer:        Communityplayers,
	}
	if err := tx.Create(&newMatch).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat pertandingan"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{
		"message": "Pertandingan telah berhasil dicatat",
		"match_id": newMatch.ID,
	})
}

func UpdateMatch(c *gin.Context) {
	matchID := c.Param("id")
	var input UpdateMatchInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := db.DB.Begin()

	var match db.Match
	if err := tx.Preload("Players").First(&match, matchID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Data pertandingan tidak ditemukan"})
		return
	}

	var cock db.Shuttlecock
	if err := tx.First(&cock, match.ShuttlecockID).Error; err == nil {
		cock.Stock += match.CocksUsed
		tx.Save(&cock)
	}

	var newCock db.Shuttlecock
	if err := tx.First(&newCock, input.ShuttlecockID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Stok merek kok baru tidak ditemukan"})
		return
	}
	if newCock.Stock < input.CocksUsed {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stok kok baru tidak mencukupi"})
		return
	}
	newCock.Stock -= input.CocksUsed
	tx.Save(&newCock)

	var newPlayers []db.Player
	tx.Where("id IN ?", input.PlayerIDs).Find(&newPlayers)
	if len(newPlayers) != 2 && len(newPlayers) != 4 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jumlah pemain harus 2 (Single) atau 4 (Double)"})
		return
	}

	if err := tx.Model(&match).Association("Players").Replace(&newPlayers); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui daftar pemain"})
		return
	}

	match.FormatScore = input.FormatScore
	match.ShuttlecockID = input.ShuttlecockID
	match.CocksUsed = input.CocksUsed
	match.DurationMinute = input.DurationMinute
	match.ScoreSet1 = input.ScoreSet1
	match.ScoreSet2 = input.ScoreSet2
	match.ScoreSet3 = input.ScoreSet3

	tx.Save(&match)
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Data pertandingan berhasil diperbarui!"})
}

func DeleteMatch(c *gin.Context) {
	matchID := c.Param("id")

	tx := db.DB.Begin()

	var match db.Match
	if err := tx.First(&match, matchID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Data pertandingan tidak ditemukan"})
		return
	}

	var cock db.Shuttlecock
	if err := tx.First(&cock, match.ShuttlecockID).Error; err == nil {
		cock.Stock += match.CocksUsed
		tx.Save(&cock)
	}

	tx.Model(&match).Association("Players").Clear()

	if err := tx.Delete(&match).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pertandingan"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Pertandingan dibatalkan, stok kok telah dikembalikan!"})
}

