package services

import (
	"context"
	"go-qchang/datasource"
	"go-qchang/models"
)

type CashierService interface {
	ChangeMoney(ctx context.Context, req models.ChangeMoneyRequest) (models.ChangeMoneyResponse, error)
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
		return res, err
	}

	changes, err := s.desk.CalculateChange(req.ChangeNeeded)
	if err != nil {
		return res, err
	}

	res.Changes = changes

	return res, nil
}
