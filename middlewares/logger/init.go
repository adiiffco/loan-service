package logger

import (
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Initialize() gin.HandlerFunc {
	return logger.SetLogger(
		logger.WithLogger(func(c *gin.Context, l zerolog.Logger) zerolog.Logger {
			return l.With().
				Str("id", requestid.Get(c)).
				Logger()
		}),
	)
}
