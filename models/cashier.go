package models

import "github.com/go-playground/validator/v10"

type (
	CashValue struct {
		Value  float64 `json:"value"`
		Amount int     `json:"amount"`
	}

	ChangeMoneyRequest struct {
		ChangeNeeded float64 `json:"change_needed" validate:"gt=0"`
	}

	ChangeMoneyResponse struct {
		Changes []CashValue `json:"changes"`
	}
)

var (
	v = validator.New()
)

func (c ChangeMoneyRequest) Validate() error {
	return v.Struct(c)
}
