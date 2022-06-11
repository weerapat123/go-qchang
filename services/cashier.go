package services

import (
	"context"
	"go-qchang/datasource"
	"go-qchang/models"

	log "github.com/sirupsen/logrus"
)

type CashierService interface {
	TransferMoneyIn(ctx context.Context, req models.TransferMoneyRequest) error
	ChangeMoney(ctx context.Context, req models.ChangeMoneyRequest) (models.ChangeMoneyResponse, error)
	TransferMoneyOut(ctx context.Context, req models.TransferMoneyRequest) error
}

type cashierService struct {
	desk datasource.CashierDesk
}

func NewCashierService(desk datasource.CashierDesk) CashierService {
	return &cashierService{desk}
}

func (s *cashierService) ChangeMoney(ctx context.Context, req models.ChangeMoneyRequest) (models.ChangeMoneyResponse, error) {
	res := models.ChangeMoneyResponse{}
	if err := req.Validate(); err != nil {
		log.Errorf("validate request failed, got error %v", err)
		return res, err
	}

	changeMoney := req.Payment - req.ProductPrice
	if changeMoney == 0 {
		return res, nil
	}

	res.ChangeMoney = changeMoney

	calculatedChanges, err := s.desk.CalculateChange(changeMoney)
	if err != nil {
		log.Errorf("calculate change failed, got error %v", err)
		return res, err
	}

	res.Changes = calculatedChanges
	return res, nil
}

func (s *cashierService) TransferMoneyIn(ctx context.Context, req models.TransferMoneyRequest) error {
	if err := req.Validate(); err != nil {
		log.Errorf("validate request failed, got error %v", err)
		return err
	}

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
		return err
	}

	for _, bc := range req.BankCoins {
		if err := s.desk.TransferMoneyOut(bc.Value, bc.Amount); err != nil {
			log.Errorf("transfer money out failed, got error %v", err)
			return err
		}
	}
	return nil
}
