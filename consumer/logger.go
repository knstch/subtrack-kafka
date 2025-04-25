package consumer

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/knstch/subtrack-libs/log"
)

type LoggerConsumerAdapter struct {
	lg Logger
}

type Logger interface {
	Error(msg string, err error, fields ...log.Message)
	Info(msg string, fields ...log.Message)
	Debug(msg string, fields ...log.Message)
	With(fields ...log.Message) *log.Logger
}

func getFields(fields map[string]interface{}) []log.Message {
	zapFields := make([]log.Message, 0, len(fields))

	for k, v := range fields {
		zapFields = append(zapFields, log.AddMessage(k, v))
	}

	return zapFields
}

func (l *LoggerConsumerAdapter) Error(msg string, err error, fields watermill.LogFields) {
	l.lg.Error(msg, err, getFields(fields)...)
}

func (l *LoggerConsumerAdapter) Info(msg string, fields watermill.LogFields) {
	l.lg.Info(msg, getFields(fields)...)
}

func (l *LoggerConsumerAdapter) Debug(msg string, fields watermill.LogFields) {
	l.lg.Debug(msg, getFields(fields)...)
}

func (l *LoggerConsumerAdapter) Trace(_ string, _ watermill.LogFields) {
	return
}

func (l *LoggerConsumerAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &LoggerConsumerAdapter{
		lg: l.lg.With(getFields(fields)...),
	}
}
