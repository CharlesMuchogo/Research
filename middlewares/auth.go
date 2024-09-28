package middlewares

import (
	"awesomeProject/auth"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			context.JSON(401, gin.H{"message": "Request does not contain an access token"})
			context.Abort()
			return
		}

		_, err := auth.ValidateToken(tokenString)
		if err != nil {

			context.JSON(401, gin.H{"message": err.Error()})
			context.Abort()
			return
		}

		context.Next()
	}
}
