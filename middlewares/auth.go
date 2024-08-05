package middlewares

import (
	"awesomeProject/auth"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			context.JSON(401, gin.H{"error": "Request does not contain an access token"})
			context.Abort()
			return
		}

		_, err := auth.ValidateToken(tokenString)
		if err != nil {

			context.JSON(401, gin.H{"error": err.Error()})
			context.Abort()
			return
		}

		context.Next()
	}
}

func HidePasswords() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a buffer to hold the response body
		var responseBody bytes.Buffer
		writer := io.MultiWriter(c.Writer, &responseBody)

		// Replace the ResponseWriter with a custom one
		c.Writer = &bodyWriter{body: writer, ResponseWriter: c.Writer}
		c.Next()

		// After the request is processed, modify the response body
		if c.Writer.Status() == 200 {
			var responseData map[string]interface{}
			if err := json.Unmarshal(responseBody.Bytes(), &responseData); err == nil {
				removePasswordsInMap(responseData)
				modifiedResponse, err := json.Marshal(responseData)
				if err == nil {
					// Write the modified response
					c.Writer.Header().Set("Content-Length", string(len(modifiedResponse)))
					c.Writer.Write(modifiedResponse)
				}
			}
		}
	}
}

type bodyWriter struct {
	gin.ResponseWriter
	body io.Writer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func removePasswordsInMap(data map[string]interface{}) {
	for k, v := range data {
		if k == "password" {
			delete(data, k)
		} else if nestedMap, ok := v.(map[string]interface{}); ok {
			removePasswordsInMap(nestedMap)
		} else if nestedArray, ok := v.([]interface{}); ok {
			for _, item := range nestedArray {
				if nestedMap, ok := item.(map[string]interface{}); ok {
					removePasswordsInMap(nestedMap)
				}
			}
		}
	}
}
