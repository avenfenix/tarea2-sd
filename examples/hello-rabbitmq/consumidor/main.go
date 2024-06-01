package main

import (
	"log"

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
	failOnError(err, "Error al establecer canal")
	defer ch.Close()
	rmq.connection = conn
	rmq.channel = ch
}

func main() {
	cargarEntorno()
	//rmq := RabbitMQ{}
	//rmq.connect()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Error al establecer conexion con RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Error al establecer canal")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil)
	failOnError(err, "Error al declarar la cola")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Fallo al registra un consumidor")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Se ha recibido un mensaje: %s\n", d.Body)
		}
	}()
	log.Println("Esperando mensajes.")
	<-forever
}
