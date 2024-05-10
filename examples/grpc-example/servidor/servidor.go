package servidor

import (
	"context"
	"log"

	pb "github.com/tarea2/grpc-example"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(c context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Peticion recibida: %v", in.GetName())
	return &pb.HelloReply{message: "Hola " + in.GetName()}, nil
}

func main() {

}
