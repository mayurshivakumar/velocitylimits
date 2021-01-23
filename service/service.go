package service

import (
	"velocitylimits/config"
	"velocitylimits/models"

	"github.com/sirupsen/logrus"
)

// TODO : It can be argued that this belongs in test, no here
//go:generate counterfeiter . Cache
type Cache interface {
	GetAccount(customerID string) *models.Account
	AddAccount(account *models.Account) *models.Account
	AddTransaction(id, customerID string)
	IsDuplicateTransaction(id, customerID string) bool
}

// Load the file.
func AttemptLoad(request *models.Request, config *config.Configurations, cache Cache) *models.Response {
	// check for duplicates
	if cache.IsDuplicateTransaction(request.ID, request.CustomerID) {
		logrus.Infoln("Ignoring duplicate txn: ", request.ID)
		return models.NewResponse(request.ID, request.CustomerID, false)
	}
	// add transactions
	cache.AddTransaction(request.ID, request.CustomerID)
	accepted := ProcessRequest(request, cache, config)
	response := models.NewResponse(request.ID, request.CustomerID, accepted)

	return response
}

// ProcessRequest ...
func ProcessRequest(request *models.Request, cache Cache, config *config.Configurations) bool {
	// Fetch the account from cache
	account := cache.GetAccount(request.CustomerID)
	// account not in cache
	if account == nil {
		account = models.NewAccount(request.CustomerID)
		account.DailyLimit = models.NewDailyLimit(request.ParsedTime, config.VelocityLimit.MaxDailyLoadLimit, config.VelocityLimit.MaxDailyTransactions)
		account.WeeklyLimit = models.NewWeeklyLimit(request.ParsedTime, config.VelocityLimit.MaxWeeklyLoadLimit)
		cache.AddAccount(account)
	} else {
		account.ResetLapsedLimits(request.ParsedTime, config.VelocityLimit.MaxDailyLoadLimit, config.VelocityLimit.MaxDailyTransactions, config.VelocityLimit.MaxWeeklyLoadLimit)
	}
	// Act on the request (if velocity limits agree)
	return account.LoadFunds(request)
}
