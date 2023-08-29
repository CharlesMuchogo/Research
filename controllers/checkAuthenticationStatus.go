package controllers
import (
	"github.com/gin-gonic/gin"
	"net/http"
)
func CheckAuthenticationStatus(context *gin.Context)  {
	context.JSON(http.StatusOK, gin.H{"message": "user session is valid"})
}