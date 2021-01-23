package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResponse(t *testing.T) {
	t.Run("returns expected response", func(t *testing.T) {
		expectedResponse := &Response{
			ID:         "1",
			CustomerID: "1",
			Accepted:   false,
		}
		actualResponse := NewResponse("1", "1", false)
		assert.Equal(t, expectedResponse, actualResponse)

	})
}
