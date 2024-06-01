package main

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func cargarEntorno() {
	err := godotenv.Load()
	failOnError(err, "Error al cargar el entorno!")
}

type RabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func (rmq *RabbitMQ) connect(url string) {
	conn, err := amqp.Dial(url)
	failOnError(err, "Error al establecer conexion con RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	defer ch.Close()
	rmq.connection = conn
	rmq.channel = ch
}

func main() {
	cargarEntorno()
	//rmq := RabbitMQ{}
	//rmq.connect("amqp://admin:1234@localhost:5672/")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Error al establecer conexion con RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil)
	failOnError(err, "Error al declarar la cola")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Fallo al publicar el mensaje")
	log.Printf(" [x] Enviar %s\n", body)
}
