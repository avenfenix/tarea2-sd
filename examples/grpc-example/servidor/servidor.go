package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	pb "github.com/tarea2/grpc_example"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedSaludadorServer
}

func (s *server) DecirHola(c context.Context, in *pb.SolicitudHola) (*pb.RespuestaHola, error) {
	log.Printf("Peticion recibida: %v", in.GetName())
	return &pb.RespuestaHola{Message: "Hola " + in.GetName()}, nil
}

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error al leer el archivo .env")
	}

	listener, err := net.Listen("tcp", ":"+os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterSaludadorServer(s, &server{})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
