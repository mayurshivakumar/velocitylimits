package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	t.Run("returns expected account", func(t *testing.T) {
		expectedAccount := &Account{
			CustomerID: "1",
		}
		actualAccount := NewAccount("1")
		assert.Equal(t, expectedAccount, actualAccount)
	})
}

func TestNewDailyLimit(t *testing.T) {
	expectedDailyLimit := &DailyLimit{
		Date:            getBeginningOfDay(time.Now()),
		MaxLoadLimit:    0,
		MaxTransactions: 0,
	}
	actualDailyLimit := NewDailyLimit(time.Now(), 0, 0)
	assert.Equal(t, expectedDailyLimit, actualDailyLimit)
}

func TestNewWeeklyLimit(t *testing.T) {
	expectedWeeklyLimit := &WeeklyLimit{
		Date:         getBeginningOfWeek(time.Now()),
		MaxLoadLimit: 0,
	}
	actualWeeklyLimit := NewWeeklyLimit(time.Now(), 0)
	assert.Equal(t, expectedWeeklyLimit, actualWeeklyLimit)
}

func TestValidateDailyLimit(t *testing.T) {
	t.Run("returns true when loading below max load limit", func(t *testing.T) {
		dailyLimit := NewDailyLimit(time.Now(), float64(2000), 3)
		valid := dailyLimit.Validate(1)
		assert.True(t, valid)
	})
	t.Run("returns true when loading exactly max load limit", func(t *testing.T) {
		dailyLimit := NewDailyLimit(time.Now(), float64(2000), 3)
		valid := dailyLimit.Validate(2000)
		assert.True(t, valid)
	})
	t.Run("returns false when loading more max load limit", func(t *testing.T) {
		dailyLimit := NewDailyLimit(time.Now(), float64(2000), 3)
		valid := dailyLimit.Validate(2001)
		assert.False(t, valid)
	})
	t.Run("returns true when loading below max transactions limit", func(t *testing.T) {
		dailyLimit := NewDailyLimit(time.Now(), float64(2000), 3)
		valid := dailyLimit.Validate(1)
		assert.True(t, valid)
	})
	t.Run("returns true when loading exactly max transactions limit", func(t *testing.T) {
		dailyLimit := NewDailyLimit(time.Now(), float64(2000), 1)
		valid := dailyLimit.Validate(2000)
		assert.True(t, valid)
	})
	t.Run("returns false when loading more max transactions limit", func(t *testing.T) {
		dailyLimit := NewDailyLimit(time.Now(), float64(2000), 0)
		valid := dailyLimit.Validate(2001)
		assert.False(t, valid)
	})
}

func TestValidateWeeklyLimit(t *testing.T) {
	t.Run("returns ture when loading below max limit", func(t *testing.T) {
		WeeklyLimit := NewWeeklyLimit(time.Now(), 20000)
		valid := WeeklyLimit.Validate(200)
		assert.True(t, valid)
	})
	t.Run("returns ture when loading equal to  max limit", func(t *testing.T) {
		WeeklyLimit := NewWeeklyLimit(time.Now(), 20000)
		valid := WeeklyLimit.Validate(20000)
		assert.True(t, valid)
	})
	t.Run("returns false when loading more than  max limit", func(t *testing.T) {
		WeeklyLimit := NewWeeklyLimit(time.Now(), 20000)
		valid := WeeklyLimit.Validate(200000)
		assert.False(t, valid)
	})
}

func TestApplyWeeklyLimit(t *testing.T) {
	t.Run("reduces weekly max", func(t *testing.T) {
		weeklyLimit := NewWeeklyLimit(time.Now(), 10)
		weeklyLimit.Apply(float64(2))
		assert.Equal(t, float64(8), weeklyLimit.MaxLoadLimit)
	})
}
func TestApplyDailyLimit(t *testing.T) {
	t.Run("reduces daily max", func(t *testing.T) {
		dailyLimit := NewDailyLimit(time.Now(), 10, 1)
		dailyLimit.Apply(float64(2))
		assert.Equal(t, float64(8), dailyLimit.MaxLoadLimit)
	})
}

func TestGetBeginningOfDay(t *testing.T) {
	now := time.Now()
	expectedDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	actualDay := getBeginningOfDay(now)
	assert.Equal(t, expectedDay, actualDay)
}

func TestGetBeginningOfWeek(t *testing.T) {
	now := time.Now()
	expectedWeek := time.Date(now.Year(), now.Month(), now.Day()+int(time.Monday-now.Weekday()), 0, 0, 0, 0, time.UTC)
	actualWeek := getBeginningOfWeek(now)
	assert.Equal(t, expectedWeek, actualWeek)
}

func TestRestLapsedLimits(t *testing.T) {
	t.Run("limits are not reset if they not before the transactions", func(t *testing.T) {
		account := NewAccount("1")
		now := time.Now()
		account.DailyLimit = NewDailyLimit(now, 1, 1)
		account.WeeklyLimit = NewWeeklyLimit(now, 1)
		account.ResetLapsedLimits(time.Now(), 2, 2, 2)
		assert.Equal(t, getBeginningOfDay(now), account.DailyLimit.Date)
		assert.Equal(t, float64(1), account.DailyLimit.MaxLoadLimit)
		assert.Equal(t, 1, account.DailyLimit.MaxTransactions)
		assert.Equal(t, getBeginningOfWeek(now), account.WeeklyLimit.Date)
		assert.Equal(t, float64(1), account.WeeklyLimit.MaxLoadLimit)

	})
	t.Run("limits are  reset if they after the transactions", func(t *testing.T) {
		account := NewAccount("1")
		yearAgo := time.Now().AddDate(-1, 0, 0)
		account.DailyLimit = NewDailyLimit(yearAgo, 1, 1)
		account.WeeklyLimit = NewWeeklyLimit(yearAgo, 1)
		now := time.Now()
		account.ResetLapsedLimits(now, 2, 2, 2)
		assert.Equal(t, getBeginningOfDay(now), account.DailyLimit.Date)
		assert.Equal(t, float64(2), account.DailyLimit.MaxLoadLimit)
		assert.Equal(t, 2, account.DailyLimit.MaxTransactions)
		assert.Equal(t, getBeginningOfWeek(now), account.WeeklyLimit.Date)
		assert.Equal(t, float64(2), account.WeeklyLimit.MaxLoadLimit)
	})
}

func TestLoadFunds(t *testing.T) {
	t.Run("returns true when loading max daily load  or weekly limit is not reached and limits are updated", func(t *testing.T) {
		account := NewAccount("528")
		account.DailyLimit = NewDailyLimit(time.Now(), 4000, 2)
		account.WeeklyLimit = NewWeeklyLimit(time.Now(), 5000)
		request, err := NewRequest("{\"id\":\"15887\",\"customer_id\":\"528\",\"load_amount\":\"$3000\",\"time\":\"2000-01-01T00:00:00Z\"}")
		require.NoError(t, err)
		success := account.LoadFunds(request)
		assert.True(t, success)
		assert.Equal(t, getBeginningOfDay(time.Now()), account.DailyLimit.Date)
		assert.Equal(t, float64(1000), account.DailyLimit.MaxLoadLimit)
		assert.Equal(t, 1, account.DailyLimit.MaxTransactions)
		assert.Equal(t, getBeginningOfWeek(time.Now()), account.WeeklyLimit.Date)
		assert.Equal(t, float64(2000), account.WeeklyLimit.MaxLoadLimit)
		assert.Equal(t, float64(3000), account.Balance)
	})
	t.Run("returns false when  when loading max daily load is reached and limits are not updated", func(t *testing.T) {
		account := NewAccount("528")
		account.DailyLimit = NewDailyLimit(time.Now(), 2000, 2)
		account.WeeklyLimit = NewWeeklyLimit(time.Now(), 5000)
		request, err := NewRequest("{\"id\":\"15887\",\"customer_id\":\"528\",\"load_amount\":\"$3000\",\"time\":\"2000-01-01T00:00:00Z\"}")
		require.NoError(t, err)
		success := account.LoadFunds(request)
		assert.False(t, success)
		assert.Equal(t, getBeginningOfDay(time.Now()), account.DailyLimit.Date)
		assert.Equal(t, float64(2000), account.DailyLimit.MaxLoadLimit)
		assert.Equal(t, 2, account.DailyLimit.MaxTransactions)
		assert.Equal(t, getBeginningOfWeek(time.Now()), account.WeeklyLimit.Date)
		assert.Equal(t, float64(5000), account.WeeklyLimit.MaxLoadLimit)
		assert.Equal(t, float64(0), account.Balance)

	})
	t.Run("returns false when  when loading max daily transactions is reached and limits are not updated", func(t *testing.T) {
		account := NewAccount("528")
		account.DailyLimit = NewDailyLimit(time.Now(), 4000, 0)
		account.WeeklyLimit = NewWeeklyLimit(time.Now(), 5000)
		request, err := NewRequest("{\"id\":\"15887\",\"customer_id\":\"528\",\"load_amount\":\"$3000\",\"time\":\"2000-01-01T00:00:00Z\"}")
		require.NoError(t, err)
		success := account.LoadFunds(request)
		assert.False(t, success)
		assert.Equal(t, getBeginningOfDay(time.Now()), account.DailyLimit.Date)
		assert.Equal(t, float64(4000), account.DailyLimit.MaxLoadLimit)
		assert.Equal(t, 0, account.DailyLimit.MaxTransactions)
		assert.Equal(t, getBeginningOfWeek(time.Now()), account.WeeklyLimit.Date)
		assert.Equal(t, float64(5000), account.WeeklyLimit.MaxLoadLimit)
		assert.Equal(t, float64(0), account.Balance)

	})
	t.Run("returns false when  when loading max weekly load is reached and limits are not updated", func(t *testing.T) {
		account := NewAccount("528")
		account.DailyLimit = NewDailyLimit(time.Now(), 4000, 1)
		account.WeeklyLimit = NewWeeklyLimit(time.Now(), 1000)
		request, err := NewRequest("{\"id\":\"15887\",\"customer_id\":\"528\",\"load_amount\":\"$3000\",\"time\":\"2000-01-01T00:00:00Z\"}")
		require.NoError(t, err)
		success := account.LoadFunds(request)
		assert.False(t, success)
		assert.Equal(t, getBeginningOfDay(time.Now()), account.DailyLimit.Date)
		assert.Equal(t, float64(4000), account.DailyLimit.MaxLoadLimit)
		assert.Equal(t, 1, account.DailyLimit.MaxTransactions)
		assert.Equal(t, getBeginningOfWeek(time.Now()), account.WeeklyLimit.Date)
		assert.Equal(t, float64(1000), account.WeeklyLimit.MaxLoadLimit)
		assert.Equal(t, float64(0), account.Balance)
	})
}
