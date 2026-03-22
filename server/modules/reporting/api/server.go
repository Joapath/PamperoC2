package api

import (
	"github.com/gin-gonic/gin"
)

func setupCORS(router *gin.Engine) {
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	setupCORS(router)
	router.GET("/health", Health)
	router.GET("/api/health", Health)
	api := router.Group("/api/v1")
	{
		reports := api.Group("/reports")
		{
			reports.POST("", GenerateReport)
			reports.GET("", ListReports)
			reports.GET("/:id", GetReport)
			reports.GET("/:id/download", DownloadReport)
			reports.DELETE("/:id", DeleteReport)
			reports.GET("/:id/ai-analysis", GetAIAnalysis)
			reports.POST("/:id/reanalyze", ReanalyzeReport)
		}

		agents := api.Group("/agents")
		{
			agents.POST("/beacon", BeaconAgent)
			agents.GET("", ListAgents)
			agents.GET("/:id", GetAgent)
			agents.POST("/:id/jobs", CreateAgentJob)
			agents.GET("/:id/jobs", ListAgentJobs)
		}

		jobs := api.Group("/jobs")
		{
			jobs.POST("/:id/results", SubmitJobResult)
		}

		api.GET("/statistics", GetStatistics)
	}
	return router
}

func StartServer(port string) error {
	router := SetupRouter()
	return router.Run(":" + port)
}
