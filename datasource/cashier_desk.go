package datasource

import (
	"go-qchang/models"
	"sync"
)

type CashierDesk interface {
	CalculateChange(change float64) ([]models.CashValue, error)
}

type desk struct {
	sync.Mutex
	BankCoins []models.CashValue
}

func New() CashierDesk {
	return &desk{
		BankCoins: []models.CashValue{
			{Value: 1000, Amount: 10},
			{Value: 500, Amount: 20},
			{Value: 100, Amount: 15},
			{Value: 50, Amount: 20},
			{Value: 20, Amount: 30},
			{Value: 10, Amount: 20},
			{Value: 5, Amount: 20},
			{Value: 1, Amount: 20},
			{Value: 0.25, Amount: 50},
		},
	}
}

func (d *desk) CalculateChange(change float64) ([]models.CashValue, error) {
	d.Lock()
	defer d.Unlock()

	return []models.CashValue{}, nil
}
