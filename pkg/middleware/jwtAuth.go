package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wangzuxianaa/tiktok-simple/pkg/token"
	"github.com/wangzuxianaa/tiktok-simple/service"
	"net/http"
)

type Token struct {
	Token string `form:"token" query:"token"`
}

func TokenChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStruct Token
		if err := c.ShouldBind(&tokenStruct); err != nil {
			c.JSON(http.StatusOK, service.Response{
				StatusCode: 1,
				StatusMsg:  "error",
			})
		}

		tokenStr := tokenStruct.Token
		if tokenStr == "" {
			c.JSON(http.StatusOK, service.Response{
				StatusCode: 1,
				StatusMsg:  "Token is not found",
			})
			c.Abort()
			return
		}
		claims, flag := token.ParseToken(tokenStr)
		if flag == false {
			c.JSON(http.StatusOK, service.Response{
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
