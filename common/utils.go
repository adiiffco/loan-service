package common

import (
	"runtime"
	"strings"

	"loanapp/models"

	"github.com/gin-gonic/gin"
)

func GetCallerMethodDetails(skip int) (resp *models.Caller) {
	pc, file, line, ok := runtime.Caller(skip)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		methodName := details.Name()
		resp = &models.Caller{
			File:   file,
			Method: methodName[strings.LastIndex(methodName, ".")+1:],
			Line:   line,
		}
	}
	return
}

func GetUserIdFromContext(c *gin.Context) int64 {
	userId, _ := c.Get("user_id")
	return userId.(int64)
}
