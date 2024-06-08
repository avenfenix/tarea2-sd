package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Message struct {
	Content    string    `json:"content"`
	ReceivedAt time.Time `json:"receivedAt"`
}

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
	conn, err := amqp.Dial("amqp://localhost:" + os.Getenv("RABBITMQ_PORT"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Crear canal
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declarar cola

	q, err := ch.QueueDeclare(
		"operations", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
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
	// Conectar a la base de datos de la maquina virtual 2
	mongoURI := "mongodb://" + os.Getenv("MONGODB_HOST") + ":" + os.Getenv("MONGODB_PORT")
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	failOnError(err, "Failed to connect to MongoDB")
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("tarea2").Collection("registros")
	// Loop

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			// Inicializar el mensaje
			var msg Message
			msg.Content = string(d.Body)
			msg.ReceivedAt = time.Now()
			// Insertar el mensaje en MongoDB
			doc := bson.D{
				{Key: "message", Value: string(d.Body)},
				{Key: "receivedAt", Value: time.Now()},
			}
			_, err := collection.InsertOne(context.TODO(), doc)
			if err != nil {
				log.Printf("Failed to insert document: %s", err)
			} else {
				log.Printf("Inserted a document: %s", d.Body)
			}
		}
	}()
	log.Printf(" Servicio de registros.")
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
