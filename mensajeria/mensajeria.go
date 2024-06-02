package main

import (
	"log"
	"os"
	"net/smtp"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func sendMail(to, subject, body string) error {
	// Correo Gmail desde el cual se mandan los correos
    from := "" 
	// Clave para el correo, se debe generar una appPassword
    password := "" 

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
	// TODO: Automatizar la obtención del correo
	to := ""
    subject := "Prueba de envío de correo en Go"
    body := "Este es un correo de prueba enviado desde un programa en Go."

    err := sendMail(to, subject, body)
    if err != nil {
        log.Fatalf("Error al enviar el correo: %s", err)
    }
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
