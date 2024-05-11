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
protoc --go_out=. --go_opt=paths=source_relative \
--go-grpc_out=. --go-grpc_opt=paths=source_relative \
proto/example.proto
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