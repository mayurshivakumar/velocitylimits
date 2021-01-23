package service_test

import (
	"testing"
	"time"

	"velocitylimits/service"
	"velocitylimits/service/servicefakes"

	"velocitylimits/models"

	"velocitylimits/cache"

	"github.com/stretchr/testify/assert"

	"velocitylimits/config"

	"github.com/stretchr/testify/require"
)

func TestAttemptLoad(t *testing.T) {
	t.Run("successful attempt to load", func(t *testing.T) {
		config := &config.Configurations{VelocityLimit: config.VelocityLimit{
			MaxDailyLoadLimit:    10,
			MaxDailyTransactions: 1,
			MaxWeeklyLoadLimit:   10,
		}}
		cache := cache.NewCache()
		request, err := models.NewRequest("{\"id\":\"15887\",\"customer_id\":\"528\",\"load_amount\":\"$3\",\"time\":\"2000-01-01T00:00:00Z\"}")
		require.NoError(t, err)
		actualResponse := service.AttemptLoad(request, config, cache)
		expectedResponse := models.NewResponse("15887", "528", true)
		assert.Equal(t, expectedResponse, actualResponse)
	})
	t.Run("returns false for duplicate request", func(t *testing.T) {
		config := &config.Configurations{VelocityLimit: config.VelocityLimit{
			MaxDailyLoadLimit:    10,
			MaxDailyTransactions: 1,
			MaxWeeklyLoadLimit:   10,
		}}
		cache := cache.NewCache()
		request, err := models.NewRequest("{\"id\":\"15887\",\"customer_id\":\"528\",\"load_amount\":\"$3\",\"time\":\"2000-01-01T00:00:00Z\"}")
		require.NoError(t, err)
		// first attempt
		actualResponse := service.AttemptLoad(request, config, cache)
		expectedResponse := models.NewResponse("15887", "528", true)
		assert.Equal(t, expectedResponse, actualResponse)
		// second attempt
		actualResponse = service.AttemptLoad(request, config, cache)
		expectedResponse = models.NewResponse("15887", "528", false)
		assert.Equal(t, expectedResponse, actualResponse)
	})
}

// this test is written a example of dependency injection.
func TestAttemptLoadWithMockedCache(t *testing.T) {
	t.Run("successful attempt to load", func(t *testing.T) {
		config := &config.Configurations{VelocityLimit: config.VelocityLimit{
			MaxDailyLoadLimit:    10,
			MaxDailyTransactions: 1,
			MaxWeeklyLoadLimit:   10,
		}}
		fakeCache := new(servicefakes.FakeCache)
		fakeCache.IsDuplicateTransactionReturns(false)
		request, err := models.NewRequest("{\"id\":\"15887\",\"customer_id\":\"528\",\"load_amount\":\"$3\",\"time\":\"2000-01-01T00:00:00Z\"}")
		require.NoError(t, err)
		actualResponse := service.AttemptLoad(request, config, fakeCache)
		expectedResponse := models.NewResponse("15887", "528", true)
		assert.Equal(t, expectedResponse, actualResponse)
	})
}

func TestProcessRequest(t *testing.T) {
	t.Run("returns true for new account", func(t *testing.T) {
		config := &config.Configurations{VelocityLimit: config.VelocityLimit{
			MaxDailyLoadLimit:    10,
			MaxDailyTransactions: 1,
			MaxWeeklyLoadLimit:   10,
		}}
		cache := cache.NewCache()
		request, err := models.NewRequest("{\"id\":\"15887\",\"customer_id\":\"528\",\"load_amount\":\"$3\",\"time\":\"2000-01-01T00:00:00Z\"}")
		require.NoError(t, err)
		actualResponse := service.ProcessRequest(request, cache, config)
		assert.True(t, actualResponse)

	})
	t.Run("returns true for existing account", func(t *testing.T) {
		config := &config.Configurations{VelocityLimit: config.VelocityLimit{
			MaxDailyLoadLimit:    10,
			MaxDailyTransactions: 1,
			MaxWeeklyLoadLimit:   10,
		}}
		cache := cache.NewCache()
		account := models.NewAccount("528")
		now := time.Now()
		account.DailyLimit = models.NewDailyLimit(now, config.VelocityLimit.MaxDailyLoadLimit, config.VelocityLimit.MaxDailyTransactions)
		account.WeeklyLimit = models.NewWeeklyLimit(now, config.VelocityLimit.MaxWeeklyLoadLimit)

		cache.AddAccount(account)

		request, err := models.NewRequest("{\"id\":\"15887\",\"customer_id\":\"528\",\"load_amount\":\"$3\",\"time\":\"2000-01-01T00:00:00Z\"}")
		require.NoError(t, err)
		actualResponse := service.ProcessRequest(request, cache, config)
		assert.True(t, actualResponse)

	})
}
