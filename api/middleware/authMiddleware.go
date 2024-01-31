package middleware

import (
	"fmt"
	"net/http"

	helper "movie-api/api/resource/user/helpers"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			msg := fmt.Sprintf("No Authorization header provider")
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": msg})
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
			c.Abort()
			return
		}

		c.Set("email_address", claims.Email_address)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("user_type", claims.User_type)
		c.Set("user_id", claims.User_id)

		c.Next()
	}
}
