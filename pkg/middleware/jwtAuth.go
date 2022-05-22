package middleware

import (
	"github.com/RaymondCode/simple-demo/controller"
	"github.com/RaymondCode/simple-demo/pkg/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TokenChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			c.JSON(http.StatusOK, controller.Response{
				StatusCode: 1,
				StatusMsg:  "Token is not found",
			})
			c.Abort()
			return
		}
		claims, flag := token.ParseToken(tokenStr)
		if flag == false {
			c.JSON(http.StatusOK, controller.Response{
				StatusCode: 1,
				StatusMsg:  "Token is not valid",
			})
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
