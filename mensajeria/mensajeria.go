package main

import (
	"context"
	"log"
	"net/smtp"
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

func sendMail(to, subject, password, from, body string) error {

	// Se Configura el servidor SMTP
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Configurar la autenticación
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Componer el mensaje
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	// Enviar el correo
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		return err
	}

	log.Println("Correo enviado correctamente")
	return nil
}

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error al leer el archivo .env")
	}
	password := os.Getenv("SENDER_PASSWORD")
	from := os.Getenv("SENDER_EMAIL")

	subject := "DiSis proteccion"
	body := "Su archivo ya se encuentra protegido y disponible en la vm2."

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
		"mensajes", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Conexión a MongoDB
	mongoURI := "mongodb://" + os.Getenv("MONGODB_HOST") + ":" + os.Getenv("MONGODB_PORT")
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	failOnError(err, "Failed to connect to MongoDB")
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("testdb").Collection("testcollection")

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
			to := string(d.Body)
			errM := sendMail(to, subject, password, from, body)

			var msg Message
			msg.Content = "Se ha enviado el correo"
			msg.ReceivedAt = time.Now()
			// Insertar el mensaje en MongoDB

			if errM != nil {
				msg.Content = "Error al enviar el correo"
			}
			doc := bson.D{
				{Key: "message", Value: msg.Content},
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
	log.Printf(" Servicio de mensajeria.")
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
