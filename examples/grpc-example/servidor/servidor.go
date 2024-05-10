package servidor

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
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(c context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Peticion recibida: %v", in.GetName())
	return &pb.HelloReply{message: "Hola " + in.GetName()}, nil
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
	pb.RegisterGreeterServer(s, &server{})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
