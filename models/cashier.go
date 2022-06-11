package models

import "github.com/go-playground/validator/v10"

type (
	CashValue struct {
		Value  float64 `json:"value"`
		Amount int     `json:"amount"`
	}

	ChangeMoneyRequest struct {
		ProductPrice float64 `json:"product_price" validate:"gt=0"`
		Payment      float64 `json:"payment" validate:"gtefield=ProductPrice"`
	}

	ChangeMoneyResponse struct {
		ChangeMoney float64     `json:"change_money"`
		Changes     []CashValue `json:"changes"`
	}

	TransferMoneyRequest struct {
		BankCoins []CashValue `json:"bank_coins"`
	}
)

var (
	v = validator.New()
)

func (c ChangeMoneyRequest) Validate() error {
	return v.Struct(c)
}
