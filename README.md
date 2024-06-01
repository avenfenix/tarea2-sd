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


## Instalar mongodb
```shell
chmod +x ./scripts/install_mongo.sh
./scripts/install_mongo.sh
```

## Instalar go
```shell
chmod +x ./scripts/install_go.sh
./scripts/install_go.sh
```

## Instalar rabbitmq-server
```shell
chmod +x ./scripts/install_rabbitmq.sh
./scripts/install_rabbitmq.sh
```

Si queremos realizar compilaciones de nuestros archivos .proto debemos instalar lo siguiente:
```shell

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


### [Ejemplo de uso gRPC y protobuffers](./examples/grpc-example/README.md)

## RabbitMQ

RabbitMQ es un intermediario de mensajes: acepta y reenvía mensajes. Puedes pensar en ello como una oficina postal: cuando depositas el correo que deseas enviar en un buzón, puedes estar seguro de que el cartero eventualmente entregará el correo a tu destinatario. En esta analogía, RabbitMQ es un buzón, una oficina postal y un cartero.

La diferencia principal entre RabbitMQ y la oficina postal es que no trabaja con papel, en lugar de eso acepta, almacena y reenvía bloques binarios de datos, es decir, mensajes.


### misc

- [Messaging Patterns](https://www.enterpriseintegrationpatterns.com/patterns/messaging/RequestReply.html)