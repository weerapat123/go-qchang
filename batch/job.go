package batch

import (
	"go-qchang/services"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func Start(svc services.CashierService) {
	c := cron.New()
	c.AddFunc("@every 5m", func() {
		defer func() {
			if r := recover(); r != nil {
				log.Warnf("recovered from system panic crash, got %v", r)
			}
		}()

		log.Infof("[Backup] Every 5 minute job")

		if err := svc.BackUpData(); err != nil {
			log.Errorf("job backup data failed, got error %v", err)
		}
	})

	c.Start()
}
