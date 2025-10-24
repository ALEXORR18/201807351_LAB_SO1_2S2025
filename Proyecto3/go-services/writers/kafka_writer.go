package writers

import (
	"context"
	"os"

	"github.com/segmentio/kafka-go"
)

type KafkaWriter struct {
	writer *kafka.Writer
}

func NewKafkaWriter() *KafkaWriter {
	brokers := os.Getenv("KAFKA_BROKERS")
	topic := os.Getenv("KAFKA_TOPIC")

	if brokers == "" {
		brokers = "localhost:9092"
	}
	if topic == "" {
		topic = "tweets-clima"
	}

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokers},
		Topic:   topic,
	})

	return &KafkaWriter{writer: w}
}

func (k *KafkaWriter) Write(msg string) error {
	return k.writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: []byte(msg),
		},
	)
}

func (k *KafkaWriter) Close() error {
	return k.writer.Close()
}
