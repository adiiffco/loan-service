package accounting

import "loanapp/models"

type Accounter interface {
	GetClient() Accounter
	FetchBalanceSheet(userId int64) []models.BalanceSheet
}
