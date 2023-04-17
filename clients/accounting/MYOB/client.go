package myob

import (
	"loanapp/clients/accounting"
	loanInterface "loanapp/clients/accounting/interface"
	"loanapp/models"
)

type MyobClient struct {
}

func (x *MyobClient) GetClient() loanInterface.Accounter {
	return accounting.NewClient("", "", nil)
}

func (x *MyobClient) FetchBalanceSheet(userId int64) []models.BalanceSheet {
	return []models.BalanceSheet{
		{
			Year:         2020,
			Month:        12,
			ProfitOrLoss: 240000,
			AssetValue:   1234,
		},
		{
			Year:         2020,
			Month:        11,
			ProfitOrLoss: 1150,
			AssetValue:   5789,
		},
		{
			Year:         2020,
			Month:        10,
			ProfitOrLoss: 2500,
			AssetValue:   22345,
		},
		{
			Year:         2020,
			Month:        9,
			ProfitOrLoss: -187000,
			AssetValue:   223452,
		},
	}
}
