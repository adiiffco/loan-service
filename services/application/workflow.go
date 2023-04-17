package application

import (
	"context"
	"loanapp/adapters/cache"
	"loanapp/adapters/logger"
	accountHelper "loanapp/clients/accounting/helper"
	"loanapp/clients/decision-engine"
	"loanapp/models"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Worker interface {
	MetaData() map[string]interface{}
	Initiate(ctx context.Context, userId int64) (models.Application, error)
	GenerateJWT(userId int64) (string, error)
	Submit(ctx context.Context, application models.ApplicationDetails) error
	FetchBalanceSheet(ctx context.Context, uuid string) (models.ApplicationDetailsRsponse, error)
	Verify(ctx context.Context, uuid string) error
	Decision(ctx context.Context, uuid string) (bool, error)
}

type Workflow struct {
	Db    DataLayer
	Cache *cache.Cache
	Log   *logger.Log
}

func InitializeWorkflow() Worker {
	return &Workflow{
		Db: InitializeDao(),
		Log: &logger.Log{
			Tag: "Application-Workflow",
		},
		Cache: cache.GetCacheInstance(),
	}
}

func (w *Workflow) GenerateJWT(userId int64) (string, error) {
	sampleSecretKey := viper.GetString("JWT_SECRET_KEY")
	token := jwt.New(jwt.SigningMethodEdDSA)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute)
	claims["user_id"] = userId

	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (w *Workflow) MetaData() map[string]interface{} {
	meta := map[string]interface{}{
		"supported_accounting_service": accountHelper.GetSupportedAccountingService(),
	}
	return meta
}

func (w *Workflow) Initiate(ctx context.Context, userId int64) (models.Application, error) {
	app, err := w.Db.Initiate(ctx, userId)
	if err != nil {
		w.Log.LogData(ctx, logrus.ErrorLevel, err.Error(), logrus.Fields{"business_user": userId})
		return models.Application{}, err
	}
	return app, nil
}

func (w *Workflow) Submit(ctx context.Context, application models.ApplicationDetails) error {
	err := w.Db.Submit(ctx, application)
	if err != nil {
		w.Log.LogData(ctx, logrus.ErrorLevel, err.Error(), logrus.Fields{"uuid": application.UUID})
	}
	return err
}

func (w *Workflow) Verify(ctx context.Context, uuid string) error {
	err := w.Db.Verify(ctx, uuid)
	if err != nil {
		w.Log.LogData(ctx, logrus.ErrorLevel, err.Error(), logrus.Fields{"uuid": uuid})
	}
	return err
}

func (w *Workflow) FetchBalanceSheet(ctx context.Context, uuid string) (models.ApplicationDetailsRsponse, error) {
	var response models.ApplicationDetailsRsponse
	applicationDetails, err := w.Db.FetchApplication(ctx, uuid)
	if err != nil {
		w.Log.LogData(ctx, logrus.ErrorLevel, err.Error(), logrus.Fields{"uuid": uuid})
		return response, err
	}
	balanceSheet, err := w.fetchBalanceSheet(ctx, applicationDetails.AccountService, applicationDetails.InitiatedBy)
	if err != nil {
		w.Log.LogData(ctx, logrus.ErrorLevel, err.Error(), logrus.Fields{"uuid": uuid})
		return response, err
	}
	response.Application = applicationDetails
	response.BalanceSheet = balanceSheet
	return response, nil
}

func (w *Workflow) fetchBalanceSheet(ctx context.Context, serviceChoice int, businessUserId int64) ([]models.BalanceSheet, error) {
	var balanceSheet []models.BalanceSheet
	choice := accountHelper.AccountingChoices(serviceChoice)
	accountService := accountHelper.GetAccountingService(choice)
	balanceSheet = (accountService.GetClient()).FetchBalanceSheet(businessUserId)
	return balanceSheet, nil
}

func (w *Workflow) Decision(ctx context.Context, uuid string) (bool, error) {
	loanApproved := false
	appDetails, err := w.FetchBalanceSheet(ctx, uuid)
	if err != nil {
		w.Log.LogData(ctx, logrus.ErrorLevel, err.Error(), logrus.Fields{"uuid": uuid})
		return loanApproved, err
	}
	preAssesmentValue := w.getPreAssesmentValue(ctx, appDetails)
	loanApproved = decision.NewClient().GetDecision(appDetails, preAssesmentValue)
	return loanApproved, nil
}

func (w *Workflow) getPreAssesmentValue(ctx context.Context, appDetails models.ApplicationDetailsRsponse) int {
	totalGain := 0.0
	totalAssetValue := 0.0
	preAssesmentValue := 20
	for _, obj := range appDetails.BalanceSheet {
		totalGain += obj.ProfitOrLoss
		totalAssetValue += obj.AssetValue
	}
	avgAssetValue := int(totalAssetValue) / len(appDetails.BalanceSheet)
	if totalGain > 0 {
		preAssesmentValue = 60
	} else if avgAssetValue > int(appDetails.LoanAmount) {
		preAssesmentValue = 100
	}
	return preAssesmentValue
}
