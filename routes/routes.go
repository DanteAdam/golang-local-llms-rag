package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	
	server.MaxMultipartMemory = 8 << 20
	server.POST("file/upload", uploadFile)
	server.DELETE("file/delete/:filename", deleteUploadedFile)
	server.POST("user/prompt", postUserPrompt)
}
