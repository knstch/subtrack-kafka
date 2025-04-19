package kafka

type KafkaTopic string

func (k KafkaTopic) String() string {
	return string(k)
}
