package batch

import (
	"go-qchang/datasource"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func New(desk datasource.CashierDesk) {
	c := cron.New()
	c.AddFunc("@every 5m", func() {
		defer func() {
			if r := recover(); r != nil {
				log.Warnf("recovered from system panic crash, got %v", r)
			}
		}()

		log.Infof("[Backup] Every 5 minute job")

		if err := desk.BackUpData(); err != nil {
			log.Errorf("job backup data failed, got error %v", err)
		}
	})

	c.Start()
}
