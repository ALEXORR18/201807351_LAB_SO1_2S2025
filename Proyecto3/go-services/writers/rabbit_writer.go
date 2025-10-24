package writers

import (
	"encoding/json"
	"os"

	"github.com/streadway/amqp"
)

// Tweet define la estructura que se enviará a RabbitMQ
type Tweet struct {
	Municipality string `json:"municipality"`
	Temperature  int    `json:"temperature"`
	Humidity     int    `json:"humidity"`
	Weather      string `json:"weather"`
}

type RabbitWriter struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue string
}

func NewRabbitWriter() (*RabbitWriter, error) {
	url := os.Getenv("RABBIT_URL")
	queue := os.Getenv("RABBIT_QUEUE")

	if url == "" {
		url = "amqp://guest:guest@localhost:5672/"
	}
	if queue == "" {
		queue = "tweets_clima"
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queue,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitWriter{
		conn:  conn,
		ch:    ch,
		queue: queue,
	}, nil
}

// WriteTweet serializa el Tweet a JSON y lo envía a RabbitMQ
func (r *RabbitWriter) WriteTweet(tweet Tweet) error {
	data, err := json.Marshal(tweet)
	if err != nil {
		return err
	}

	return r.ch.Publish(
		"",
		r.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
}

// Close cierra el canal y la conexión
func (r *RabbitWriter) Close() {
	r.ch.Close()
	r.conn.Close()
}
