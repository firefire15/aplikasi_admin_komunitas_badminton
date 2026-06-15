package controllers

import (
	"net/http"
	"aplikasi_admin_komunitas_badminton/db"
	"github.com/gin-gonic/gin"
)

type FinancialReportResponse struct {
	TotalExpenses  float64 `json:"total_expenses"`  
	TotalIncomes   float64 `json:"total_incomes"`   
	NetBalance     float64 `json:"net_balance"`      
}

type CockReportResponse struct {
	Brand        string `json:"brand"`
	CurrentStock int    `json:"current_stock_balls"` 
	TotalUsed    int    `json:"total_used_balls"`   
	TotalReturned int   `json:"total_returned_balls"`
}

func GetFinancialReport(c *gin.Context) {
	communityID := c.MustGet("community_id").(uint)

	var report FinancialReportResponse

	db.DB.Model(&db.Receipt{}).
		Where("community_id = ?", communityID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&report.TotalExpenses)

	db.DB.Model(&db.Payment{}).
		Where("community_id = ? AND is_confirmed = ?", communityID, true).
		Select("COALESCE(SUM(total_paid), 0)").
		Scan(&report.TotalIncomes)

	report.NetBalance = report.TotalIncomes - report.TotalExpenses

	c.JSON(http.StatusOK, gin.H{
		"status":           "success",
		"financial_report": report,
	})
}

// 2. LAPORAN LOGISTIK KOK (Stok, Penggunaan, & Retur)
func GetShuttlecockReport(c *gin.Context) {
	communityID := c.MustGet("community_id").(uint)

	var inventories []db.Shuttlecock
	db.DB.Where("community_id = ?", communityID).Find(&inventories)

	var reportDetails []CockReportResponse

	for _, inv := range inventories {
		var totalUsed int64

		db.DB.Model(&db.Match{}).
			Where("shuttlecock_id = ?", inv.ID).
			Select("COALESCE(SUM(cocks_used), 0)").
			Scan(&totalUsed)

		reportDetails = append(reportDetails, CockReportResponse{
			Brand:         inv.Brand,
			CurrentStock:  inv.Stock,
			TotalUsed:     int(totalUsed),
			TotalReturned: inv.Returned,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":             "success",
		"shuttlecock_report": reportDetails,
	})
}