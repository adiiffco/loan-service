package models

type ApplicationDetails struct {
	UUID            string  `json:"uuid" binding:"uuid"`
	BusinessName    string  `json:"name" binding:"required"`
	YearEstablished int     `json:"year" binding:"required"`
	AccountService  int     `json:"service" binding:"required"`
	LoanAmount      float64 `json:"loan_amount" binding:"required"`
}
