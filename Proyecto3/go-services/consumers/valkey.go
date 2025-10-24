package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

// Contexto global
var ctx = context.Background()

// Tweet con los mismos campos que genera Locust
type Tweet struct {
	Municipio string `json:"Municipality"`
	Temp      int    `json:"Temperature"`
	Humedad   int    `json:"Humidity"`
	Clima     string `json:"Weather"`
}

// Conexión a Valkey (Redis)
func getValkeyClient() *redis.Client {
	addr := os.Getenv("VALKEY_SERVICE_URL")
	if addr == "" {
		addr = "valkey-service:6379"
	}
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
}

// Consumidor de RabbitMQ → guarda en Valkey
func startRabbitConsumer(rdb *redis.Client) {
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@rabbitmq-service.clima-app.svc.cluster.local:5672/"
	}

	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		log.Fatal("❌ Error conectando a RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("❌ Error abriendo canal:", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("tweets_clima", true, false, false, false, nil)
	if err != nil {
		log.Fatal("❌ Error declarando cola:", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal("❌ Error al consumir mensajes:", err)
	}

	fmt.Println("🐇 Esperando mensajes desde RabbitMQ...")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var tweet Tweet
			if err := json.Unmarshal(d.Body, &tweet); err != nil {
				log.Printf("❌ Error deserializando mensaje: %v", err)
				log.Printf("%s", string(d.Body))
				continue
			}

			msgBody := d.Body // d es tu mensaje de RabbitMQ
			log.Println("Mensaje raw recibido:", string(msgBody))

			data, _ := json.Marshal(tweet)
			err := rdb.LPush(ctx, "clima_tweets", data).Err()
			if err != nil {
				log.Printf("❌ Error guardando en Valkey: %v", err)
			} else {
				fmt.Println(tweet)
				fmt.Println("✅ Tweet guardado en Valkey:", string(data))
			}
		}
	}()

	<-forever
}

func main() {
	for {
		rdb := getValkeyClient()
		ping, err := rdb.Ping(ctx).Result()
		if err != nil {
			fmt.Println("❌ No se pudo conectar a Valkey, reintentando en 3s...")
			time.Sleep(3 * time.Second)
			continue
		}
		fmt.Println("✅ Conexión a Valkey establecida:", ping)
		startRabbitConsumer(rdb)
	}
}
