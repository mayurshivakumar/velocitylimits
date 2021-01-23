package models

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Account...
type Account struct {
	CustomerID  string
	Balance     float64
	DailyLimit  *DailyLimit
	WeeklyLimit *WeeklyLimit
}

// DailyLimit...
type DailyLimit struct {
	Date            time.Time
	MaxLoadLimit    float64
	MaxTransactions int
}

// WeeklyLimit...
type WeeklyLimit struct {
	Date         time.Time
	MaxLoadLimit float64
}

// NewDailyLimit...
func NewDailyLimit(d time.Time, maxLoadLimit float64, maxTransactions int) *DailyLimit {
	return &DailyLimit{
		Date:            getBeginningOfDay(d),
		MaxLoadLimit:    maxLoadLimit,
		MaxTransactions: maxTransactions,
	}
}

// NewWeeklyLimit...
func NewWeeklyLimit(d time.Time, maxLoadLimit float64) *WeeklyLimit {
	return &WeeklyLimit{
		Date:         getBeginningOfWeek(d),
		MaxLoadLimit: maxLoadLimit,
	}
}

// NewAccount...
func NewAccount(customerID string) *Account {
	return &Account{
		CustomerID: customerID,
	}
}

// Validate Daily Limit...
func (dl *DailyLimit) Validate(amount float64) bool {
	if dl.MaxLoadLimit-amount < 0 {
		return false
	}
	if dl.MaxTransactions-1 < 0 {
		return false
	}
	return true
}

// Apply DailyLimit
func (dl *DailyLimit) Apply(amount float64) {
	dl.MaxLoadLimit -= amount
	dl.MaxTransactions--
}

// Validate Weekly limit
func (wl *WeeklyLimit) Validate(amount float64) bool {
	if wl.MaxLoadLimit-amount < 0 {
		return false
	}
	return true
}

// Apply  weekly limit
func (wl *WeeklyLimit) Apply(amount float64) {
	wl.MaxLoadLimit -= amount
}

// ResetLapsedLimits ...
func (a *Account) ResetLapsedLimits(t time.Time, maxDailyLoadLimit float64, maxTransactions int, maxWeeklyLoadLimit float64) {
	transactionDay := getBeginningOfDay(t)
	if transactionDay.After(a.DailyLimit.Date) {
		a.DailyLimit.Date = transactionDay
		a.DailyLimit.MaxLoadLimit = maxDailyLoadLimit
		a.DailyLimit.MaxTransactions = maxTransactions
	}
	transactionWeek := getBeginningOfWeek(t)
	if transactionWeek.After(a.WeeklyLimit.Date) {
		a.WeeklyLimit.Date = transactionWeek
		a.WeeklyLimit.MaxLoadLimit = maxWeeklyLoadLimit
	}
}

// LoadFunds ...
func (a *Account) LoadFunds(r *Request) bool {
	// Validate if daily limits
	if a.DailyLimit.Validate(r.ParsedAmount) == false {
		logrus.Debugln("Daily limit reached. request rejected: ", r.ID)
		return false
	}
	// Validate if weekly limits
	if a.WeeklyLimit.Validate(r.ParsedAmount) == false {
		logrus.Debugln("Weekly limit reached. request rejected: ", r.ID)
		return false
	}
	a.Balance += r.ParsedAmount
	// Update the limits after acting on this transactions
	a.DailyLimit.Apply(r.ParsedAmount)
	a.WeeklyLimit.Apply(r.ParsedAmount)
	logrus.Debugln("Transaction approved: ", r.ID)
	return true
}

// getBeginningOfDay
func getBeginningOfDay(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
}

// getBeginningOfWeek
func getBeginningOfWeek(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day()+int(time.Monday-d.Weekday()), 0, 0, 0, 0, time.UTC)
}
