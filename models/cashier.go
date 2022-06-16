package models

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	FieldProductPrice string = "ProductPrice"
	FieldPayment      string = "Payment"
	FieldBankCoins    string = "BankCoins"
	FieldValue        string = "Value"
	FieldAmount       string = "Amount"
)

type (
	CashValue struct {
		Value  float64 `json:"value" validate:"required,value"`
		Amount int     `json:"amount" validate:"gt=0"`
	}

	ChangeMoneyRequest struct {
		ProductPrice float64     `json:"product_price" validate:"gt=0"`
		Payment      float64     `json:"payment" validate:"gt=0,gtefield=ProductPrice"`
		BankCoins    []CashValue `json:"bank_coins" validate:"min=1,dive"`
	}

	ChangeMoneyResponse struct {
		ChangeMoney float64     `json:"change_money"`
		Changes     []CashValue `json:"changes"`
	}

	TransferMoneyRequest struct {
		BankCoins []CashValue `json:"bank_coins" validate:"min=1,dive"`
	}

	CheckResponse struct {
		Total     float64     `json:"total"`
		BankCoins []CashValue `json:"bank_coins"`
	}
)

var (
	v = validator.New()
)

func init() {
	v.RegisterValidation("value", func(fl validator.FieldLevel) bool {
		value := fl.Field().Float()
		switch value {
		case 1000, 500, 100, 50, 20, 10, 5, 1, 0.25:
			return true
		}
		return false
	})
}

func (req ChangeMoneyRequest) Validate() error {
	if err := v.Struct(req); err != nil {
		return err
	}
	total := 0.0
	for _, bc := range req.BankCoins {
		total += bc.Value * float64(bc.Amount)
	}
	if total != req.Payment {
		return fmt.Errorf("Payment and back/coin is inconsistent.")
	}
	return nil
}

func (req TransferMoneyRequest) Validate() error {
	return v.Struct(req)
}

func TransformErrorMessage(err error) error {
	switch err.(type) {
	case validator.ValidationErrors:
		errs := err.(validator.ValidationErrors)
		msgs := make([]string, 0, len(errs))

		for _, e := range errs {
			switch e.Field() {
			case FieldProductPrice:
				msgs = append(msgs, "Must provide product price.")
			case FieldPayment:
				if e.Tag() == "gtefield" {
					msgs = append(msgs, "Payment is not enough.")
				} else {
					msgs = append(msgs, "Must provide payment.")
				}
			case FieldBankCoins:
				msgs = append(msgs, "Must provide bank/coin.")
			case FieldValue:
				if e.Tag() == "value" {
					msgs = append(msgs, "Bank/coin value is invalid.")
				} else {
					msgs = append(msgs, "Must provide bank/coin value.")
				}
			case FieldAmount:
				msgs = append(msgs, "Must provide bank/coin amount.")
			}
		}

		return fmt.Errorf(strings.Join(msgs, " "))
	default:
		return err
	}
}
