/* package main

import (
	"context"
	"log"
	"net"

	pb "go-services/common/proto" // Ajusta el path a tu proto generado

	"google.golang.org/grpc"
)

// Servidor gRPC que implementa WeatherTweetService
type weatherServer struct {
	pb.UnimplementedWeatherTweetServiceServer
}

// Implementación del método SendTweet
func (s *weatherServer) SendTweet(ctx context.Context, req *pb.WeatherTweetRequest) (*pb.WeatherTweetResponse, error) {
	// Aquí puedes procesar el tweet, guardarlo o hacer log
	log.Printf("Tweet recibido: municipio=%v, temperatura=%d, humedad=%d, clima=%v",
		req.Municipality, req.Temperature, req.Humidity, req.Weather)

	// Respuesta al cliente gRPC
	return &pb.WeatherTweetResponse{
		Status: "Tweet recibido correctamente ✅",
	}, nil
}

func main() {
	// Escuchar conexiones en el puerto 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("No se pudo abrir el puerto: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterWeatherTweetServiceServer(grpcServer, &weatherServer{})

	log.Println("Servidor gRPC corriendo en el puerto 50051 🚀")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error al iniciar el servidor gRPC: %v", err)
	}
}
*/

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "go-services/common/proto"
	"go-services/writers"

	"google.golang.org/grpc"
)

type weatherServer struct {
	pb.UnimplementedWeatherTweetServiceServer
	kafkaWriter  *writers.KafkaWriter
	rabbitWriter *writers.RabbitWriter
}

func (s *weatherServer) SendTweet(ctx context.Context, req *pb.WeatherTweetRequest) (*pb.WeatherTweetResponse, error) {
	// Mapear enums a strings
	municipios := map[pb.Municipalities]string{
		pb.Municipalities_mixco:     "mixco",
		pb.Municipalities_guatemala: "guatemala",
		pb.Municipalities_amatitlan: "amatitlan",
		pb.Municipalities_chinautla: "chinautla",
	}

	climas := map[pb.Weathers]string{
		pb.Weathers_sunny:  "sunny",
		pb.Weathers_cloudy: "cloudy",
		pb.Weathers_rainy:  "rainy",
		pb.Weathers_foggy:  "foggy",
	}

	msg := fmt.Sprintf("Municipio=%v, Temp=%d, Humedad=%d, Clima=%v",
		municipios[req.Municipality],
		req.Temperature,
		req.Humidity,
		climas[req.Weather])

	msg2 := writers.Tweet{
		Municipality: municipios[req.Municipality],
		Temperature:  int(req.Temperature),
		Humidity:     int(req.Humidity),
		Weather:      climas[req.Weather],
	}

	log.Println("Tweet recibido:", msg)

	// Enviar a Kafka
	if err := s.kafkaWriter.Write(msg); err != nil {
		log.Println("Error enviando a Kafka:", err)
	}

	// Enviar a RabbitMQ
	if err := s.rabbitWriter.WriteTweet(msg2); err != nil {
		log.Println("Error enviando a RabbitMQ:", err)
	}

	return &pb.WeatherTweetResponse{
		Status: "Tweet recibido y enviado a Kafka y RabbitMQ ✅",
	}, nil
}

func main() {
	// Inicializar writers usando variables de entorno
	kafkaWriter := writers.NewKafkaWriter()

	rabbitWriter, err := writers.NewRabbitWriter()
	if err != nil {
		log.Fatalf("Error inicializando RabbitMQ: %v", err)
	}

	defer kafkaWriter.Close()
	defer rabbitWriter.Close()

	// Inicializar servidor gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("No se pudo abrir el puerto: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterWeatherTweetServiceServer(grpcServer, &weatherServer{
		kafkaWriter:  kafkaWriter,
		rabbitWriter: rabbitWriter,
	})

	log.Println("Servidor gRPC corriendo en el puerto 50051 🚀")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error al iniciar el servidor gRPC: %v", err)
	}
}
