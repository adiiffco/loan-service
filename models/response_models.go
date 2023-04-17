package models

type BalanceSheet struct {
	Year         int     `json:"year"`
	Month        int     `json:"month"`
	ProfitOrLoss float64 `json:"profitOrLoss"`
	AssetValue   float64 `json:"assetsValue"`
}

type ApplicationDetailsRsponse struct {
	Application
	BalanceSheet []BalanceSheet `json:"balance_sheet"`
}
