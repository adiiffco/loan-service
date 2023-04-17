package accounting

import (
	"fmt"

	"loanapp/adapters/logger"
	"loanapp/clients"
	loanInterface "loanapp/clients/accounting/interface"
	"loanapp/models"
)

type BaseClient struct {
	ApiClient *clients.ApiClient
	Log       *logger.Log
	Headers   map[string]string
}

func NewClient(baseURL, clientTag string, baseHeaders map[string]string) loanInterface.Accounter {
	apiClient := clients.NewClient()
	apiClient.BaseURL = baseURL
	return &BaseClient{
		ApiClient: apiClient,
		Log: &logger.Log{
			Tag: fmt.Sprintf("AccountingClient-%s", clientTag),
		},
		Headers: baseHeaders,
	}
}

func (b *BaseClient) GetClient() loanInterface.Accounter {
	return nil
}

func (b *BaseClient) FetchBalanceSheet(userId int64) []models.BalanceSheet {
	return nil
}
