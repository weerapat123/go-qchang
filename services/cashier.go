package services

import (
	"context"
	"go-qchang/datasource"
	"go-qchang/models"

	log "github.com/sirupsen/logrus"
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
	res := models.ChangeMoneyResponse{ChangeMoney: make([]models.CashValue, 0)}
	if err := req.Validate(); err != nil {
		log.Errorf("validate request failed, got error: %v", err)
		return res, err
	}

	changeMoney := req.Cash - req.ProductPrice
	if changeMoney == 0 {
		return res, nil
	}

	calculatedChangeMoney, err := s.desk.CalculateChange(changeMoney)
	if err != nil {
		log.Errorf("calculate change failed, got error: %v", err)
		return res, err
	}

	res.ChangeMoney = calculatedChangeMoney

	return res, nil
}
