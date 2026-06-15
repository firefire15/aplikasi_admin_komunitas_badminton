package routes

import(
	"aplikasi_admin_komunitas_badminton/controllers"
	"github.com/gin-gonic/gin"
	"aplikasi_admin_komunitas_badminton/helper"
)

func StartAPIServer() *gin.Engine {

	router := gin.Default()
	badminton := router.Group("/api/badminton")
	{
		badminton.POST("/admin/register", controllers.RegisterAdmin)
		badminton.POST("/admin/login", controllers.LoginAdmin)
		badminton.POST("/communities/register", controllers.CreateCommunity)
		badminton.GET("/communities", controllers.GetCommunities)
		badminton.POST("/players", controllers.CreatePlayer) 

		protected := badminton.Group("/")
		protected.Use(helper.JWTMiddleware())
		{       
			protected.GET("/communities/players/:id", controllers.GetMyCommunityPlayers)
			protected.POST("/communities/players", controllers.AddCommunityPlayers)

			protected.POST("/court/", controllers.AddCourt)
			protected.GET("/court/:id", controllers.GetCourt)

			protected.POST("/schedules/book", controllers.BookCourtAndSchedule)

			protected.POST("/shuttlecocks/buy", controllers.BuyShuttlecock)
			protected.GET("/shuttlecocks/stock/:id", controllers.GetShuttlecockStock)
			protected.POST("/shuttlecocks/return", controllers.ReturnShuttlecock)

			protected.POST("/matches", controllers.RecordMatch)            
			protected.PUT("/matches/:id", controllers.UpdateMatch)         
			protected.DELETE("/matches/:id", controllers.DeleteMatch)      
		
			protected.GET("/reports/financial", controllers.GetFinancialReport) 
			protected.GET("/reports/shuttlecocks", controllers.GetShuttlecockReport) 

			protected.POST("/billing/generate", controllers.GenerateSessionBilling)
			protected.PUT("/billing/confirm/:id", controllers.ConfirmPaymentOK) 

		}
	}
	return router
}