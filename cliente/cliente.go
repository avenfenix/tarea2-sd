package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	pb "github.com/tarea2/proteccion/proto"
)

func main() {

	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error al leer el archivo .env")
	}

	if len(os.Args) < 4 {
		fmt.Println("Usa: go run cliente <rut> <correo> <ruta>")
		os.Exit(1)
	}

	rut := os.Args[1]
	correo := os.Args[2]
	ruta := os.Args[3]

	// Manejo de archivo

	// Abrir el archivo desde la ruta especificada
	archivo, err := os.Open(ruta)
	if err != nil {
		log.Fatalf("No se pudo abrir el archivo: %v", err)
	}
	defer archivo.Close()

	// Leer los datos del archivo
	info, err := archivo.Stat()
	if err != nil {
		log.Fatalf("No se pudo leer la información del archivo: %v", err)
	}
	tamaño := info.Size()
	datos := make([]byte, tamaño)
	_, err = archivo.Read(datos)
	if err != nil {
		log.Fatalf("No se pudo leer el contenido del archivo: %v", err)
	}

	connection, err := grpc.Dial(os.Getenv("GRPC_HOST")+":"+os.Getenv("GRPC_PORT"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer connection.Close()

	cliente := pb.NewProtectorClient(connection)

	solicitud := &pb.SolicitudProteger{
		Rut:    rut,
		Correo: correo,
		File:   datos,
	}

	respuesta, err := cliente.Proteger(context.Background(), solicitud)
	if err != nil {
		log.Fatalf("No se pudo proteger: %v", err)
	}

	log.Printf("Respuesta: %s", respuesta.Message)

}
