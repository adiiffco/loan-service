package decision

import (
	"loanapp/adapters/logger"
	"loanapp/clients"
	"loanapp/models"

	"github.com/spf13/viper"
)

type BaseClient struct {
	ApiClient *clients.ApiClient
	Log       *logger.Log
	Headers   map[string]string
}

func NewClient() *BaseClient {
	apiClient := clients.NewClient()
	apiClient.BaseURL = viper.GetString("DECISION_URL")
	return &BaseClient{
		ApiClient: apiClient,
		Log: &logger.Log{
			Tag: "BaseClient-Decision",
		},
		Headers: map[string]string{},
	}
}

func (b *BaseClient) GetDecision(appDetails models.ApplicationDetailsRsponse, preAssesmentValue int) bool {
	return preAssesmentValue > 20
}
