package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	pb "github.com/tarea2/grpc_example"
)

func main() {

	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error al leer el archivo .env")
	}

	connection, err := grpc.Dial(os.Getenv("GRPC_HOST")+":"+os.Getenv("GRPC_PORT"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}

	defer connection.Close()

	cliente := pb.NewSaludadorClient(connection)

	name := "Mundo"

	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	respuesta, err := cliente.DecirHola(context.Background(), &pb.SolicitudHola{Name: name})
	if err != nil {
		log.Fatalf("No se pudo saludar: %v", err)
	}
	log.Printf("Saludo: %s", respuesta.Message)
}
