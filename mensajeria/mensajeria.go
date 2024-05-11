package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error al leer el archivo .env")
	}
	// Conexion
	conn, err := amqp.Dial("amqp://guest:guest@" + os.Getenv("RABBITMQ_HOST") + ":" + os.Getenv("RABBITMQ_PORT"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Crear canal
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declarar cola

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Podemos leer los mensajes desde el canal

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	// Loop

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf(" Servicio de mensajeria.")
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
