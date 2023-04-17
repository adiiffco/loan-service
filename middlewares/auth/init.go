package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func BasicAuth() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{viper.GetString("BASIC_USERNAME"): viper.GetString("BASIC_PASSWORD")})
}

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		headers := c.Request.Header
		authToken := headers.Get("authorization")
		if authToken != "" {
			token, err := validateJWT(authToken)
			if err != nil || !token.Valid {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			claims := token.Claims.(jwt.MapClaims)
			userId := claims["user_id"].(int64)
			newContext := context.WithValue(c.Request.Context(), "user_id", userId)
			c.Request = c.Request.WithContext(newContext)
			c.Next()
			return
		}
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func validateJWT(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		var err error
		t, ok := token.Method.(*jwt.SigningMethodECDSA)
		if !ok {
			err = errors.New("not authorized")
			t = nil
		}
		return t, err
	})
}
