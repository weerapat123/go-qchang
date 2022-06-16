package datasource

import (
	"encoding/csv"
	"fmt"
	"go-qchang/models"
	"io"
	"math"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type CashierDesk interface {
	TransferMoneyIn(value float64, amount int) error
	TransferMoneyOut(value float64, amount int) error
	CalculateChange(change float64) ([]models.CashValue, error)
	GetTotalBankCoin() []models.CashValue
	ResetBankCoin()
	BackUpData() error
}

const DefaultDataPath string = "assets/bankcoins.csv"

type desk struct {
	dataPath  string
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

func New(path string) CashierDesk {
	d := &desk{dataPath: path}
	if err := d.loadData(); err != nil {
		panic(err)
	}

	if !d.isDataValid() {
		panic("data is corrupted")
	}
	return d
}

func (d *desk) loadData() error {
	bankCoinsDefault := make([]models.CashValue, 0, len(valueMapToIndex))

	file, err := os.Open(d.dataPath)
	if err != nil {
		return err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	// Header
	_, err = csvReader.Read()
	if err != nil {
		return err
	}

	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if len(record) == 2 {
			value, err := strconv.ParseFloat(record[0], 64)
			if err != nil {
				return err
			}
			amount, err := strconv.Atoi(record[1])
			if err != nil {
				return err
			}
			bankCoinsDefault = append(bankCoinsDefault, models.CashValue{value, amount})
		}
	}
	d.BankCoins = bankCoinsDefault
	return nil
}

func (d *desk) isDataValid() bool {
	defer func() {
		// Data is invalid and some row is missing that will lead to pacic
		if r := recover(); r != nil {
			log.Warnf("data is not valid, got %v", r)
		}
	}()

	for key, pointer := range valueMapToIndex {
		if d.BankCoins[pointer.index].Value != key {
			log.Warnf("data is not valid, got data value %v and key %v", d.BankCoins[pointer.index].Value, key)
			return false
		}
	}
	return true
}

func (d *desk) CalculateChange(change float64) ([]models.CashValue, error) {
	tmpBankCoins := make([]models.CashValue, len(d.BankCoins))
	copy(tmpBankCoins, d.BankCoins)

	changes := make([]models.CashValue, 0, len(valueMapToIndex))
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
		return nil, fmt.Errorf("No change available")
	}

	d.BankCoins = tmpBankCoins
	return changes, nil
}

func (d *desk) TransferMoneyIn(value float64, amount int) error {
	pointer, ok := valueMapToIndex[value]
	if !ok {
		return fmt.Errorf("Invalid value")
	}

	if (d.BankCoins[pointer.index].Amount + amount) > pointer.limit {
		return fmt.Errorf("Bank/coin of %v is at limit", value)
	}

	d.BankCoins[pointer.index].Amount += amount
	return nil
}

func (d *desk) TransferMoneyOut(value float64, amount int) error {
	pointer, ok := valueMapToIndex[value]
	if !ok {
		return fmt.Errorf("Invalid value")
	}

	if (d.BankCoins[pointer.index].Amount - amount) < 0 {
		return fmt.Errorf("Bank/coin of %v is not enough", value)
	}

	d.BankCoins[pointer.index].Amount -= amount
	return nil
}

func (d *desk) GetTotalBankCoin() []models.CashValue {
	tmpBankCoins := make([]models.CashValue, len(d.BankCoins))
	copy(tmpBankCoins, d.BankCoins)

	return tmpBankCoins
}

func (d *desk) ResetBankCoin() {
	for i := range d.BankCoins {
		d.BankCoins[i].Amount = 0
	}
}

func (d *desk) BackUpData() error {
	file, err := os.Create(d.dataPath)
	if err != nil {
		return err
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
