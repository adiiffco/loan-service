package application

import (
	"loanapp/common"
	"loanapp/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Application struct {
	W Worker
}

func InitController() *Application {
	return &Application{
		W: InitializeWorkflow(),
	}
}

func (a *Application) External(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"healthCheck": "External",
	})
}

func (a *Application) Internal(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"healthCheck": "Internal",
	})
}

func (a *Application) Authorize(c *gin.Context) {
	testUserID := viper.GetInt64("USER_ID_TEST")
	token, err := a.W.GenerateJWT(testUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	c.Header("token", token)
	c.JSON(http.StatusOK, gin.H{})
}

func (a *Application) MetaData(c *gin.Context) {
	meta := a.W.MetaData()
	c.JSON(http.StatusOK, meta)
}

func (a *Application) Initiate(c *gin.Context) {
	userId := common.GetUserIdFromContext(c)
	application, err := a.W.Initiate(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	c.JSON(http.StatusOK, application)
}

func (a *Application) SubmitDetails(c *gin.Context) {
	var request models.ApplicationDetails
	err := a.W.Submit(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (a *Application) VerifyDetails(c *gin.Context) {
	uuid := c.Query("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	err := a.W.Verify(c.Request.Context(), uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (a *Application) BalanceSheet(c *gin.Context) {
	uuid := c.Query("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	response, err := a.W.FetchBalanceSheet(c.Request.Context(), uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (a *Application) Decision(c *gin.Context) {
	uuid := c.Query("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	loanApproved, _ := a.W.Decision(c.Request.Context(), uuid)
	c.JSON(http.StatusOK, loanApproved)
}
