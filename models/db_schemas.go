package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApplicationStatus int

const (
	INITIATED ApplicationStatus = iota
	SUBMITTED
	VERIFIED
)

type Application struct {
	Id              int64             `gorm:"AUTO_INCREMENT"`
	UUID            string            `gorm:"column:uuid;default:not null" json:"uuid"`
	BusinessName    string            `gorm:"column:business_name;default:null" json:"business_name"`
	YearEstablished int               `gorm:"column:year;default:null" json:"year"`
	LoanAmount      float64           `gorm:"column:loan_amount;default:null" json:"loan_amount"`
	AccountService  int               `gorm:"column:account_service;default:null" json:"account_service"`
	Status          ApplicationStatus `gorm:"column:status;default:not null" json:"status"`
	InitiatedBy     int64             `gorm:"column:initiated_by;default:not null" json:"initiated_by"`
	CreatedAt       int64             `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       int64             `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Application) TableName() string {
	return "application"
}

func (a *Application) BeforeCreate(tx *gorm.DB) (err error) {
	a.UUID = uuid.New().String()
	return
}
