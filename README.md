# Tarea 2
Desarrollo de tarea 2 de sistemas distribuidos.

El inicio de esta tarea comienza con el final de la primera en el cual terminamos con un sistema de proteccion que funciona pero no puede atender todas las peticiones.

Para comenzar a solucionar este problema dividimos nuestra aplicacion en 3 componentes:

* Cliente (VM 1)
* Servicio de proteccion (VM 2)
* Servicio de registros y mensajeria (VM 3)

Nuestra solucion debera ser capaz de:
 
* Realizar la protección de archivos PDF ya desarrollada
* Guardar  registros  de  las  operaciones  y  sus  resultados  una  vez  hecho  el  procesamiento del archivo. 
* Notificar al cliente para el cual el archivo está siendo protegido, a través de un mensaje por correo electrónico.


# Ejecucion

## Cliente

Correr en la primera maquina virtual (VM 1)

```shell
go run ./cliente <rut> <correo> <ruta>
```

## Servicio proteccion

Correr en la segunda maquina virtual (VM 2)
```shell
go run ./proteccion
```

## Servicio registros
Correr en la tercera maquina virtual (VM 3)


```shell
go run ./registros
```

## Servicio mensajeria
Correr en la tercera maquina virtual (VM 3)

```shell
go run ./mensajeria
```


# Instalacion y configuracion

Todas las maquinas
```shell
# GOLANG
wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz && sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz && echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
```

Instalacion VM 2 (jammy)
```shell
# MongoDB
sudo apt update
sudo apt install gnupg wget apt-transport-https ca-certificates software-properties-common
wget -qO- \
  https://pgp.mongodb.com/server-7.0.asc | \
  gpg --dearmor | \
  sudo tee /usr/share/keyrings/mongodb-server-7.0.gpg >/dev/null

echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-7.0.gpg ] \
  https://repo.mongodb.org/apt/ubuntu $(lsb_release -cs)/mongodb-org/7.0 multiverse" | \
  sudo tee -a /etc/apt/sources.list.d/mongodb-org-7.0.list
sudo apt update
sudo apt install mongodb-org

# Comandos Mongo
sudo systemctl start mongod
sudo systemctl daemon-reload
sudo systemctl status mongod
sudo systemctl stop mongod
sudo systemctl restart mongod
netstat -plntu
mongosh

# gRPC
sudo apt install -y protobuf-compiler
protoc --version  # Ensure compiler version is 3+

# plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
```


# Contenido

## Interface Description Language (IDL)

Describe una interfaz en un lenguaje neutral, que permite la comunicación entre componentes de software desarrollados en diferentes lenguajes de programación, como por ejemplo entre componentes escritos en C++ y otros escritas en Java.

### Usos

* Son utilizadas con frecuencia en el software de las llamadas a procedimiento remoto (RPC, Remote Procedure Call), lo que permite a los sistemas de computadoras utilizar lenguajes y sistemas operativos diferentes. IDL ofrece un puente entre dos sistemas diferentes.

## Protobuffers

Es como JSON, excepto que es más pequeño y rápido, y genera enlaces nativos de lenguaje. Defines cómo quieres que se estructuren tus datos una vez, luego puedes usar código fuente generado especial para escribir y leer fácilmente tus datos estructurados desde y hacia una variedad de flujos de datos y utilizando una variedad de lenguajes.

Los protocol buffers son una combinación del lenguaje de definición (creado en archivos .proto), el código que el compilador de protos genera para interactuar con los datos, las bibliotecas de tiempo de ejecución específicas del lenguaje, el formato de serialización para los datos que se escriben en un archivo (o se envían a través de una conexión de red), y los datos serializados.

![alt text](https://protobuf.dev/images/protocol-buffers-concepts.png)


* [Lo basico de protobuffer en Go](https://protobuf.dev/getting-started/gotutorial/)

## gRPC

En gRPC, una aplicación cliente puede llamar directamente a un método en una aplicación servidor en una máquina diferente como si fuera un objeto local, lo que facilita la creación de aplicaciones y servicios distribuidos. Como en muchos sistemas de RPC, gRPC se basa en la idea de definir un servicio, especificando los métodos que pueden ser llamados de forma remota con sus parámetros y tipos de retorno. En el lado del servidor, el servidor implementa esta interfaz y ejecuta un servidor gRPC para manejar las llamadas del cliente. En el lado del cliente, el cliente tiene un stub (llamado simplemente cliente en algunos lenguajes) que proporciona los mismos métodos que el servidor.

![alt text](https://grpc.io/img/landing-2.svg)

* [Introduction to gRPC](https://grpc.io/docs/what-is-grpc/introduction/)
* [gRPC explicado - Youtube](https://www.youtube.com/watch?v=NHw2cjcMN9g&t=60s)
* [Tutorial basico gRPC en Go](https://grpc.io/docs/languages/go/basics/)


# Desarrollo

Para entender como funciona gRPC y tambien los protobuffers usaremos solo las dos primeras maquinas virtuales. En la VM 1 tendremos al cliente y en la VM2 al servidor. Aunque por ahora el unico servicio sera decir Hola. 

- examples/
  - grpc-example/
    - .env
    - cliente/
      - cliente.go
      - go.mod
    - proto/
      - hello.proto
    - servidor/
      - servidor.go
      - go.mod

.env file

```env
VM2=
GRPC_HOST=${VM2}
GRPC_PORT=8080
```

hello.proto
```go
// Protocol Buffer versión 3
syntax = "proto3";

// Definición del paquete para organizar los mensajes y servicios
package grpc_example;

// Opción para especificar el paquete de Go que se utilizará en el código generado
option go_package = "github.com/tarea2/grpc_example";

// Definición del servicio "Greeter"
service Greeter {
    // Método "SayHello" que recibe un mensaje de tipo HelloRequest y devuelve un mensaje de tipo HelloResponse
    rpc SayHello (HelloRequest) returns (HelloResponse) {}
}

// Definición del mensaje HelloRequest
message HelloRequest {
    string name = 1;
}

// Definición del mensaje HelloResponse
message HelloResponse {
    string message = 1;
}
```

Compiling .proto file

```shell
protoc --go_out=. --go_opt=paths=source_relative \
--go-grpc_out=. --go-grpc_opt=paths=source_relative \
proto/example.proto
```

paquetes

```shell
go get -u github.com/joho/godotenv
```