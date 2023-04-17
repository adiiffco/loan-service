package application

import (
	"context"
	"loanapp/adapters/logger"
	"loanapp/adapters/mysql"
	"loanapp/models"

	"gorm.io/gorm"
)

type DataLayer interface {
	Initiate(ctx context.Context, userId int64) (models.Application, error)
	Submit(ctx context.Context, application models.ApplicationDetails) error
	Verify(ctx context.Context, uuid string) error
	FetchApplication(ctx context.Context, uuid string) (models.Application, error)
}

type Dao struct {
	ORM *gorm.DB
	Log *logger.Log
}

func InitializeDao() DataLayer {
	return &Dao{
		ORM: mysql.GetDbInstance(),
		Log: &logger.Log{
			Tag: "Application-Dao",
		},
	}
}

func (d *Dao) Initiate(ctx context.Context, userId int64) (models.Application, error) {
	application := models.Application{
		InitiatedBy: userId,
		Status:      models.INITIATED,
	}
	result := d.ORM.WithContext(ctx).FirstOrCreate(&application)
	return application, result.Error
}

func (d *Dao) Submit(ctx context.Context, application models.ApplicationDetails) error {
	result := d.ORM.WithContext(ctx).Model(models.Application{}).Where("uuid = ?", application.UUID).Updates(map[string]interface{}{
		"business_name":   application.BusinessName,
		"year":            application.YearEstablished,
		"account_service": application.AccountService,
		"loan_amount":     application.LoanAmount,
		"status":          models.SUBMITTED,
	})
	return result.Error
}

func (d *Dao) Verify(ctx context.Context, uuid string) error {
	result := d.ORM.WithContext(ctx).Model(models.Application{}).Where("uuid = ? and status = ?", uuid, models.SUBMITTED).Updates(map[string]interface{}{
		"status": models.VERIFIED,
	})
	return result.Error
}

func (d *Dao) FetchApplication(ctx context.Context, uuid string) (models.Application, error) {
	application := models.Application{}
	result := d.ORM.WithContext(ctx).Where("uuid", uuid).First(&application)
	return application, result.Error
}
