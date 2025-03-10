package kafka

import (
	"github.com/IBM/sarama"
	"log/slog"
)

type Producer struct {
	logger   *slog.Logger
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaProducer(logger *slog.Logger, brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		logger:   logger,
		producer: producer,
		topic:    topic,
	}, nil
}

func (p *Producer) SendMessage(key, value string) error {
	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}

	partition, offset, err := p.producer.SendMessage(message)
	if err != nil {
		return err
	}

	p.logger.Info("Сообщение отправлено: partition=%d, offset=%d, key=%s, value=%s", partition, offset, key, value)
	return nil

}

// Close закрывает продюсер.
func (p *Producer) Close() error {
	return p.producer.Close()
}
