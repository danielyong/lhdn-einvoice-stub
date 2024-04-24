package main

import (
	"lhdn-dummy/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	handlers.Init()

	r := gin.Default()

	r.POST("/connect/token", handlers.IntermediaryLogin)

	r.GET("/api/v1.0/documenttypes", handlers.ValidateAccessToken, handlers.GetDocumentTypes)
  r.GET("/api/v1.0/documenttypes/:document_type_id", handlers.ValidateAccessToken, handlers.GetDocumentType)

  r.GET("/api/v1.0/documents/recent", handlers.ValidateAccessToken, handlers.ReallySelectStarSubmission)
  r.GET("/api/v1.0/documents/search", handlers.ValidateAccessToken, handlers.ReallySelectStarSubmission)
	r.GET("/api/v1.0/documents/:document_id/raw", handlers.ValidateAccessToken, handlers.GetDocument)
	r.GET("/api/v1.0/documents/:document_id/details", handlers.ValidateAccessToken, handlers.GetDocument)
	r.POST("/api/v1.0/documentsubmissions", handlers.ValidateAccessToken, handlers.SubmitDocument)
	r.GET("/api/v1.0/documentsubmissions/:submission_id", handlers.ValidateAccessToken, handlers.SelectStarSubmission)
	r.PUT("/api/v1.0/documents/state/:document_id/state", handlers.ValidateAccessToken, handlers.UpdateDocument)

	r.Run()
}
