package cliente

import (
	"fmt"

	pb "github.com/tarea2/grpc-example"
)

func main() {
	mensaje := pb.HelloRequest{
		name: "Mundo",
	}

	fmt.Println("Mensaje: ", mensaje)

}
