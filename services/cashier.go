package services

import (
	"context"
	"fmt"
	"go-qchang/datasource"
	"go-qchang/models"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

type CashierService interface {
	TransferMoneyIn(ctx context.Context, req models.TransferMoneyRequest) error
	ChangeMoney(ctx context.Context, req models.ChangeMoneyRequest) (models.ChangeMoneyResponse, bool, error)
	TransferMoneyOut(ctx context.Context, req models.TransferMoneyRequest) error
	Check(ctx context.Context) models.CheckResponse
}

type cashierService struct {
	sync.Mutex
	desk datasource.CashierDesk
}

func NewCashierService(desk datasource.CashierDesk) CashierService {
	return &cashierService{desk: desk}
}

func (s *cashierService) ChangeMoney(ctx context.Context, req models.ChangeMoneyRequest) (models.ChangeMoneyResponse, bool, error) {
	res := models.ChangeMoneyResponse{}
	if err := req.Validate(); err != nil {
		log.Errorf("validate request failed, got error %v", err)
		return res, false, models.TransformErrorMessage(err)
	}

	changeMoney := req.Payment - req.ProductPrice
	if changeMoney == 0 {
		return res, false, nil
	}

	res.ChangeMoney = changeMoney

	s.Lock()
	defer s.Unlock()

	calculatedChanges, err := s.desk.CalculateChange(changeMoney)
	if err != nil {
		log.Errorf("calculate change failed, got error %v", err)
		return res, false, fmt.Errorf("%s. Please fill in bank/coins.", err)
	}

	res.Changes = calculatedChanges

	// Put payment money into cashier desk
	errs := make([]string, 0, len(req.BankCoins))
	for _, bc := range req.BankCoins {
		err := s.desk.TransferMoneyIn(bc.Value, bc.Amount)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) != 0 {
		return res, true, fmt.Errorf("%s. Please take some of it out.", strings.Join(errs, ". "))
	}

	return res, true, nil
}

func (s *cashierService) TransferMoneyIn(ctx context.Context, req models.TransferMoneyRequest) error {
	if err := req.Validate(); err != nil {
		log.Errorf("validate request failed, got error %v", err)
		return models.TransformErrorMessage(err)
	}

	s.Lock()
	defer s.Unlock()

	for _, bc := range req.BankCoins {
		if err := s.desk.TransferMoneyIn(bc.Value, bc.Amount); err != nil {
			log.Errorf("transfer money in failed, got error %v", err)
			return err
		}
	}
	return nil
}

func (s *cashierService) TransferMoneyOut(ctx context.Context, req models.TransferMoneyRequest) error {
	if err := req.Validate(); err != nil {
		log.Errorf("validate request failed, got error %v", err)
		return models.TransformErrorMessage(err)
	}

	s.Lock()
	defer s.Unlock()

	for _, bc := range req.BankCoins {
		if err := s.desk.TransferMoneyOut(bc.Value, bc.Amount); err != nil {
			log.Errorf("transfer money out failed, got error %v", err)
			return err
		}
	}
	return nil
}

func (s *cashierService) Check(ctx context.Context) models.CheckResponse {
	total := 0.0
	bankcoins := s.desk.GetTotalBankCoin()

	for _, bc := range bankcoins {
		total += bc.Value * float64(bc.Amount)
	}

	return models.CheckResponse{
		Total:     total,
		BankCoins: bankcoins,
	}
}
