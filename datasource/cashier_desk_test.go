package datasource

import (
	"go-qchang/models"
	"testing"
)

func Test_desk_CalculateChange(t *testing.T) {
	desk := New()

	changeMoney, err := desk.CalculateChange(1778.25)
	if err != nil {
		t.Errorf("desk.CalculateChange() got error %v", err)
	}

	expected := []models.CashValue{
		{1000, 1},
		{500, 1},
		{100, 2},
		{50, 1},
		{20, 1},
		{5, 1},
		{1, 3},
		{0.25, 1},
	}

	if len(changeMoney) != len(expected) {
		t.Errorf("desk.CalculateChange() = %v, want %v", changeMoney, expected)
	}

	for i := range changeMoney {
		if changeMoney[i].Value != expected[i].Value || changeMoney[i].Amount != expected[i].Amount {
			t.Errorf("desk.CalculateChange() = %v, want %v", changeMoney, expected)
		}
	}
}
