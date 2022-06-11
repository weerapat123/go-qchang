package datasource

import (
	"fmt"
	"go-qchang/models"
	"math"
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

	tmpBankCoins := make([]models.CashValue, len(d.BankCoins))
	copy(tmpBankCoins, d.BankCoins)

	changes := make([]models.CashValue, 0, 9)
	remaining := change

	for i, bc := range tmpBankCoins {
		if remaining <= 0 {
			break
		}

		calculatedAmount := math.Floor(remaining / bc.Value)

		if calculatedAmount != 0 && bc.Amount != 0 {
			if int(calculatedAmount) <= bc.Amount {
				remaining -= bc.Value * calculatedAmount
				changes = append(changes, models.CashValue{Value: bc.Value, Amount: int(calculatedAmount)})
				tmpBankCoins[i].Amount -= int(calculatedAmount)
			} else {
				remaining -= bc.Value * float64(bc.Amount)
				changes = append(changes, models.CashValue{Value: bc.Value, Amount: bc.Amount})
				tmpBankCoins[i].Amount = 0
			}
		}
	}

	if remaining != 0 {
		return nil, fmt.Errorf("no change available")
	}

	d.BankCoins = tmpBankCoins
	return changes, nil
}
