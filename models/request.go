package models

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Request ...
type Request struct {
	ID           string    `json:"id"`
	CustomerID   string    `json:"customer_id"`
	Amount       string    `json:"load_amount"`
	Time         string    `json:"time"`
	ParsedAmount float64   `json:"-"`
	ParsedTime   time.Time `json:"-"`
}

// NewRequest ...
func NewRequest(reqStr string) (*Request, error) {
	var r Request
	var err error

	if err = json.Unmarshal([]byte(reqStr), &r); err != nil {
		logrus.Errorln("Error parsing line: ", err)
		return nil, err
	}

	if r.ParsedAmount, err = strconv.ParseFloat(strings.Trim(r.Amount, "$"), 64); err != nil {
		logrus.Errorln("Error parsing amount: ", err)
		return nil, err
	}

	if r.ParsedTime, err = time.Parse(time.RFC3339, r.Time); err != nil {
		logrus.Errorln("Error parsing time: ", err)
		return nil, err
	}

	return &r, nil
}
