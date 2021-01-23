package models

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRequest(t *testing.T) {
	t.Run("returns expected response", func(t *testing.T) {
		parsedAmount, _ := strconv.ParseFloat(strings.Trim("$100", "$"), 64)
		parsedTime, _ := time.Parse(time.RFC3339, "2000-01-01T06:08:12Z")
		expectedRequest := &Request{
			ID:           "1",
			CustomerID:   "1",
			Amount:       "$100",
			Time:         "2000-01-01T06:08:12Z",
			ParsedAmount: parsedAmount,
			ParsedTime:   parsedTime,
		}
		actualRequest, err := NewRequest("{\"id\":\"1\",\"customer_id\":\"1\",\"load_amount\":\"$100\",\"time\":\"2000-01-01T06:08:12Z\"}")
		require.NoError(t, err)
		assert.Equal(t, expectedRequest, actualRequest)
	})
	t.Run("returns error when convert to json is failed", func(t *testing.T) {
		_, err := NewRequest("\"id\":\"1\",\"customer_id\":\"1\",\"load_amount\":\"$100\",\"time\":\"2000-01-01T06:08:12Z\"}")
		require.Error(t, err)
	})
	t.Run("returns error when parsing invalid amount string", func(t *testing.T) {
		_, err := NewRequest("{\"id\":\"1\",\"customer_id\":\"1\",\"load_amount\":\"@100\",\"time\":\"2000-01-01T06:08:12Z\"}")
		require.Error(t, err)
	})
	t.Run("returns error when parsing invalid time string", func(t *testing.T) {
		_, err := NewRequest("{\"id\":\"1\",\"customer_id\":\"1\",\"load_amount\":\"$100\",\"time\":\"2000-0101T06:08:12Z\"}")
		require.Error(t, err)
	})
}
