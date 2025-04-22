package outbox

import (
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/knstch/subtrack-kafka/producer"
	kafkaPkg "github.com/knstch/subtrack-kafka/topics"

	"go.uber.org/zap"
)

type OutboxListener struct {
	cron     *cron.Cron
	producer *producer.Producer
	db       *gorm.DB
}

func NewOutboxListener(kafkaAddr string, dbDsn string, lg *zap.Logger) (*cron.Cron, error) {
	db, err := gorm.Open(postgres.Open(dbDsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	cronProducer := producer.NewProducer(kafkaAddr)

	c := cron.New()
	if _, err = c.AddFunc("@every 10s", func() {
		var outbox []Outbox
		if err = db.Model(&Outbox{}).Where("sent_at IS NULL").Find(&outbox).Error; err != nil {
			lg.Error("error getting outbox from database", zap.Error(err))
			return
		}

		lg.Info("got items from outbox", zap.Any("amount", len(outbox)))

		for i := range outbox {
			if err = cronProducer.SendMessage(kafkaPkg.KafkaTopic(outbox[i].Topic), outbox[i].Key, outbox[i].Payload); err != nil {
				lg.Error("error sending message to kafka", zap.Error(err), zap.Any("id", outbox[i].ID))
				continue
			}

			if err = db.Model(&Outbox{}).Where("id = ?", outbox[i].ID).Update("sent_at", time.Now()).Error; err != nil {
				lg.Error("error updating outbox from database", zap.Error(err), zap.Any("id", outbox[i].ID))
				break
			}
		}

		lg.Info("cycle is done!✨")
	}); err != nil {
		return nil, err
	}

	return c, nil
}
