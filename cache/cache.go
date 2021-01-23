package cache

import (
	"velocitylimits/models"
)

// Cache ...
type Cache struct {
	accounts     map[string]*models.Account
	transactions map[string]struct{}
}

// NewCache ...
func NewCache() *Cache {
	return &Cache{
		accounts:     make(map[string]*models.Account),
		transactions: make(map[string]struct{}),
	}
}

// GetAccount ...
func (s *Cache) GetAccount(customerID string) *models.Account {
	if acc, ok := s.accounts[customerID]; ok {
		return acc
	}
	return nil
}

// AddAccountToStore ...
func (s *Cache) AddAccount(account *models.Account) *models.Account {
	s.accounts[account.CustomerID] = account
	return account
}

// AddTransaction ...
func (s *Cache) AddTransaction(id, customerID string) {
	s.transactions[id+customerID] = struct{}{}
}

// IsDuplicateTransaction ...
func (s *Cache) IsDuplicateTransaction(id, customerID string) bool {
	if _, ok := s.transactions[id+customerID]; ok {
		return true
	}
	return false
}
