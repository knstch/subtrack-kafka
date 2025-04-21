package outbox

import "time"

const (
	OutboxMigrationUp = `
	CREATE TABLE IF NOT EXISTS outbox (
	    id SERIAL PRIMARY KEY,
	    topic VARCHAR(255) NOT NULL,
	    key VARCHAR(255) NOT NULL,
	    payload JSONB NOT NULL,
	    sent_at TIMESTAMP NULL,
	    created_at TIMESTAMP DEFAULT NOW()
	)
`
	OutboxMigrationDown = `
	DROP TABLE IF EXISTS outbox;
`
)

type Outbox struct {
	ID        uint
	Topic     string
	Key       string
	Payload   []byte
	SentAt    *time.Time
	CreatedAt time.Time
}

func (Outbox) TableName() string {
	return "outbox"
}
