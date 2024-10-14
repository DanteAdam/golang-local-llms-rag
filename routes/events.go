package routes

import (
	"fmt"
	"golang-llms/utils"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type PromptType struct {
	Prompt  string `json:"prompt"`
	Message string `json:"message"`
}

func uploadFile(context *gin.Context) {
	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not upload file",
		})
		return
	}

	log.Println(file.Filename)

	sanitizedFileName := filepath.Base(file.Filename)
	err = context.SaveUploadedFile(file, "storage/"+sanitizedFileName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not save file",
		})
		return
	}
	context.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

func deleteUploadedFile(context *gin.Context) {
	fileName := context.Param("filename")

	err := utils.DeleteFile(fileName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not delete file",
		})
		return
	}
	context.String(http.StatusOK, fmt.Sprintf("'%s' deleted succesfully!", fileName))

}

func postUserPrompt(context *gin.Context) {
	var prompt PromptType
	err := context.BindJSON(&prompt)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid JSON format",
		})
		return
	}
	store := utils.Setup()
	utils.SetUserPrompt(prompt.Prompt)
	resDocs := utils.Search(store)
	response, err := utils.Execute(resDocs)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate response",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Received event with prompt: %s", prompt.Prompt)
	context.JSON(http.StatusOK, gin.H{
		"message":  "Send prompt and receive response succesfully!",
		"prompt":   prompt.Prompt,
		"response": response,
	})
}
