package cache

import (
	"testing"

	"velocitylimits/models"

	"github.com/stretchr/testify/assert"
)

func TestNewCache(t *testing.T) {
	t.Run("returns expected cache", func(t *testing.T) {
		expectedCache := &Cache{
			accounts:     make(map[string]*models.Account),
			transactions: make(map[string]struct{}),
		}
		actualCache := NewCache()
		assert.Equal(t, expectedCache, actualCache)
	})
}

func TestGetAccount(t *testing.T) {
	t.Run("returns expected account", func(t *testing.T) {
		cache := NewCache()
		expectedAccount := &models.Account{CustomerID: "1"}
		cache.accounts[expectedAccount.CustomerID] = expectedAccount
		actualAccount := cache.GetAccount(expectedAccount.CustomerID)
		assert.Equal(t, expectedAccount, actualAccount)
	})
}

func TestAddAccount(t *testing.T) {
	t.Run("adds account", func(t *testing.T) {
		cache := NewCache()
		expectedAccount := &models.Account{CustomerID: "1"}
		cache.AddAccount(expectedAccount)
		actualAccount := cache.accounts[expectedAccount.CustomerID]
		assert.Equal(t, expectedAccount, actualAccount)
	})
	t.Run("returns nil when no account not found", func(t *testing.T) {
		cache := NewCache()
		expectedAccount := &models.Account{CustomerID: "1"}
		cache.AddAccount(expectedAccount)
		actualAccount := cache.accounts["2"]
		assert.Nil(t, actualAccount)
	})
}

func TestAddTransaction(t *testing.T) {
	t.Run("adds transaction", func(t *testing.T) {
		cache := NewCache()
		cache.AddTransaction("11", "2")
		duplicate := cache.transactions["11"+"2"]
		assert.NotNil(t, duplicate)
	})
}

func TestCacheIsDuplicateTransaction(t *testing.T) {
	t.Run("returns true when there is a duplicate transaction", func(t *testing.T) {
		cache := NewCache()
		cache.AddTransaction("11", "2")
		duplicate := cache.IsDuplicateTransaction("11", "2")
		assert.True(t, duplicate)

	})
	t.Run("returns false when there is no transaction", func(t *testing.T) {
		cache := NewCache()
		cache.AddTransaction("11", "2")
		duplicate := cache.IsDuplicateTransaction("111", "21")
		assert.False(t, duplicate)
	})
}
