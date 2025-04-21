package outbox

import (
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/knstch/subtrack-kafka/producer"
	kafkaPkg "github.com/knstch/subtrack-kafka/topics"
)

type OutboxCronWorker struct {
	cron     *cron.Cron
	producer *producer.Producer
	db       *gorm.DB
}

func NewOutboxCronWorker(addr string, dsn string) (*cron.Cron, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	cronProducer := producer.NewProducer(addr)

	c := cron.New()
	if _, err = c.AddFunc("@every 10s", func() {
		var outbox []Outbox
		if err = db.Model(&Outbox{}).Where("sent_at IS NULL").Find(&outbox).Error; err != nil {
			return
		}

		for i := range outbox {
			if err = cronProducer.SendMessage(kafkaPkg.KafkaTopic(outbox[i].Topic), outbox[i].Key, outbox[i].Payload); err != nil {
				continue
			}

			if err = db.Model(&Outbox{}).Where("id = ?", outbox[i].ID).Update("sent_at", time.Now()).Error; err != nil {
				break
			}
		}
	}); err != nil {
		return nil, err
	}

	return c, nil
}
