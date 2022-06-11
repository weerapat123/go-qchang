package datasource

import (
	"encoding/csv"
	"fmt"
	"go-qchang/models"
	"io"
	"math"
	"os"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
)

type CashierDesk interface {
	FillMoneyIn(value float64, amount int) error
	CalculateChange(change float64) ([]models.CashValue, error)
	TakeMoneyOut(value float64, amount int) error
	BackUpData() error
}

const dataPath string = "assets/bankcoins.csv"

type desk struct {
	sync.Mutex
	BankCoins []models.CashValue
}

type cashPointer struct {
	index int
	limit int
}

var valueMapToIndex = map[float64]cashPointer{
	1000: {0, 10},
	500:  {1, 20},
	100:  {2, 15},
	50:   {3, 20},
	20:   {4, 30},
	10:   {5, 20},
	5:    {6, 20},
	1:    {7, 20},
	0.25: {8, 50},
}

func New() CashierDesk {
	tmp := getBankCoins()
	if !isDataValid(tmp) {
		panic("data is corrupted")
	}
	return &desk{
		BankCoins: tmp,
	}
}

func getBankCoins() []models.CashValue {
	bankCoinsDefault := make([]models.CashValue, 0, 9)

	file, err := os.Open(dataPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	// Header
	csvReader.Read()

	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		value, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			panic(err)
		}
		amount, err := strconv.Atoi(record[1])
		if err != nil {
			panic(err)
		}
		bankCoinsDefault = append(bankCoinsDefault, models.CashValue{value, amount})
	}
	return bankCoinsDefault
}

func isDataValid(data []models.CashValue) bool {
	defer func() {
		// Data is invalid and some row is missing that will lead to pacic
		if r := recover(); r != nil {
			log.Warnf("data is not valid, got %v", r)
		}
	}()

	for key, pointer := range valueMapToIndex {
		if data[pointer.index].Value != key {
			log.Warnf("data is not valid, got data value %v and key %v", data[pointer.index].Value, key)
			return false
		}
	}
	return true
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

func (d *desk) FillMoneyIn(value float64, amount int) error {
	index := valueMapToIndex[value].index
	if (d.BankCoins[index].Amount + amount) > valueMapToIndex[value].limit {
		return fmt.Errorf("amount bank/coin of %v is at limit", value)
	}

	d.BankCoins[index].Amount += amount
	return nil
}

func (d *desk) TakeMoneyOut(value float64, amount int) error {
	index := valueMapToIndex[value].index
	if (d.BankCoins[index].Amount - amount) < 0 {
		return fmt.Errorf("current amount of bank/coin of %v is not enough", value)
	}
	d.BankCoins[index].Amount -= amount
	return nil
}

func (d *desk) BackUpData() error {
	file, err := os.OpenFile(dataPath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"value", "amount"}); err != nil {
		return fmt.Errorf("write header in backup file failed, error: %w", err)
	}

	for _, bc := range d.BankCoins {
		if err := writer.Write([]string{fmt.Sprint(bc.Value), fmt.Sprint(bc.Amount)}); err != nil {
			return fmt.Errorf("backup data failed, error: %w", err)
		}
	}

	return nil
}
