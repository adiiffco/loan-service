package request

import (
	"context"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func AddRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := requestid.Get(c)
		newContext := context.WithValue(c.Request.Context(), "request_id", rid)
		c.Request = c.Request.WithContext(newContext)
		c.Next()
	}
}
