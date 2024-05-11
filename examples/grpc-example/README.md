# Ejemplo de gRPC y Protobuffers

Para entender como funciona gRPC y tambien los protobuffers usaremos solo las dos primeras maquinas virtuales. En la VM 1 tendremos al cliente y en la VM2 al servidor. 

## Interfaz de comunicacion

Utilizamos un archivo .proto que definira nuestra interfaz de comunicacion. Compilaremos el archivo .proto con `protoc` y usaremos el codigo generado en el **cliente** y **servidor** de nuestra aplicacion de ejemplo.


### Definicion
Definicion de la version de protobuf a utilizar.

```proto
syntax = "proto3";
```

Definimos el paquete para organizar los mensajes y servicios. Tambien especificamos el nombre del paquete de Go que se utilizará en el código generado

```proto
package grpc_example;
option go_package = "github.com/tarea2/grpc_example";
```

El servicio que entregara el servidor gRPC sera de saludar. Tendremos un servicio llamado `Saludador` y una rutina rpc llamada `DecirHola`. La rutina recibe un mensaje del tipo `SolicitudHola` y envia un mensaje de respuesta del tipo `RespuestaHola` 

```proto
service Saludador {
    rpc DecirHola (SolicitudHola) returns (RespuestaHola) {}
}

message SolicitudHola {
    string name = 1;
}

message RespuestaHola {
    string message = 1;
}
```

### Compilando archivo .proto

```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/example.proto
```

### Usando el codigo generado

Para usar el codigo lo tratamos como un modulo de Go

Configurando modulo

```shell
cd examples/grpc_example/proto/
go mod init github.com/tarea2/grpc_example
```

Agregar en cliente.go y servidor.go
```go
import (
	pb "github.com/tarea2/grpc_example"
)
```

### Codigo servidor

```go
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
```

### Codigo cliente

```go
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
```
### Corriendo ejemplo

```shell
# VM 2
go run ./examples/grpc-example/servidor

# VM 1
go run ./examples/grpc-example/cliente <nombre>
```

- grpc_example/
  - cliente/
    - cliente.go
    - go.mod
  - proto/
    - example.proto
    - codigo_generado_grpc.pb.go
    - codigo_generado.pb.go
    - go.mod
  - servidor/
    - servidor.go
    - go.mod